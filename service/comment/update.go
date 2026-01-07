package comment

import (
	"base_project/domain"
	"base_project/request"
	"base_project/response"
	sharedErrs "base_project/shared/errors"
	"base_project/util/constant"
	"base_project/util/logger"
	"context"
	"time"
)

func (b *base) Update(ctx context.Context, input *request.UpdateComment) (*response.Comment, error) {
	logger.Info(ctx, "Update Comment with request: %v", input)

	authUser := constant.AuthUserFromCtx(ctx)
	if authUser == nil {
		return nil, sharedErrs.UnauthorizedErr
	}

	comment, err := b.repository.FindCommentByIDAndUserID(ctx, input.ID, authUser.ID)
	if err != nil {
		return nil, err
	}

	result, err := b.repository.UpdateComment(ctx, &domain.Comment{
		BaseModel: domain.BaseModel{
			ID:        comment.ID,
			UpdatedAt: time.Now(),
		},
		UserID:  authUser.ID,
		Comment: input.Comment,
		Rating:  input.Rating,
	})
	if err != nil {
		return nil, err
	}

	return response.NewCommentFromDomain(result), nil
}
