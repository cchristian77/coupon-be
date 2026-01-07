package post

import (
	"base_project/domain"
	"base_project/request"
	"base_project/response"
	sharedErrs "base_project/shared/errors"
	"base_project/util/constant"
	"base_project/util/logger"
	"context"
	"errors"
	"time"

	"gorm.io/gorm"
)

func (b *base) Update(ctx context.Context, input *request.UpsertPost) (*response.Post, error) {
	logger.Info(ctx, "Update Post with request: %v", input)

	authUser := constant.AuthUserFromCtx(ctx)
	if authUser == nil {
		return nil, sharedErrs.UnauthorizedErr
	}

	post, err := b.repository.FindPostByIDAndUserID(ctx, input.ID, authUser.ID)
	if err != nil {
		return nil, err
	}

	if post.IsPublished() {
		return nil, sharedErrs.NewBusinessValidationErr("Published post cannot be updated. Please draft post before updating.")
	}

	slugExists, err := b.repository.FindPostBySlug(ctx, input.Slug)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}
	if slugExists != nil && slugExists.ID != post.ID {
		return nil, sharedErrs.NewBusinessValidationErr("Update Failed. Post with slug '%s' already exists.", post.Slug)
	}

	result, err := b.repository.UpdatePost(ctx, &domain.Post{
		BaseModel: domain.BaseModel{
			ID:        post.ID,
			UpdatedAt: time.Now(),
		},
		UserID: authUser.ID,
		Slug:   toSlug(input.Slug),
		Title:  input.Title,
		Body:   input.Body,
	})
	if err != nil {
		return nil, err
	}

	return response.NewPostFromDomain(result), nil
}
