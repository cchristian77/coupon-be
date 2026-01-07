package post

import (
	"base_project/domain"
	"base_project/domain/enums"
	"base_project/request"
	"base_project/response"
	sharedErrs "base_project/shared/errors"
	"base_project/util/constant"
	"base_project/util/logger"
	"context"
	"errors"
	"regexp"
	"strings"
	"time"

	"gorm.io/gorm"
)

func (b *base) Store(ctx context.Context, input *request.UpsertPost) (*response.Post, error) {
	logger.Info(ctx, "Store Post with request: %v", input)

	authUser := constant.AuthUserFromCtx(ctx)
	if authUser == nil {
		return nil, sharedErrs.UnauthorizedErr
	}

	slugExists, err := b.repository.FindPostBySlug(ctx, input.Slug)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}
	if slugExists != nil {
		return nil, sharedErrs.NewBusinessValidationErr(
			"Create Failed. Post with slug '%s' already exists.", input.Slug)
	}

	now := time.Now()
	post, err := b.repository.CreatePost(ctx, &domain.Post{
		BaseModel: domain.BaseModel{
			CreatedAt: now,
			UpdatedAt: now,
		},
		UserID: authUser.ID,
		Slug:   toSlug(input.Slug),
		Title:  input.Title,
		Body:   input.Body,
		Status: enums.DRAFTPostStatus,
	})
	if err != nil {
		return nil, err
	}

	return response.NewPostFromDomain(post), nil
}

func toSlug(input string) string {
	slug := strings.ToLower(input)

	// Replace any character that is NOT a-z or 0-9 with a hyphen
	// The regex pattern `[^a-z0-9]+` matches one or more non-alphanumeric characters
	reg := regexp.MustCompile(`[^a-z0-9]+`)
	slug = reg.ReplaceAllString(slug, "-")

	// rim hyphens from the start and end
	slug = strings.Trim(slug, "-")

	return slug
}
