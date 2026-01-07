package comment

import (
	"context"
	"coupon_be/domain"
	"coupon_be/request"
	"coupon_be/response"
	sharedErrs "coupon_be/shared/errors"
	"coupon_be/util/constant"
	"coupon_be/util/logger"
	"errors"
	"time"

	"gorm.io/gorm"
)

func (b *base) Store(ctx context.Context, input *request.CreateComment) (*response.Comment, error) {
	logger.Info(ctx, "Store Comment with request: %v", input)

	authUser := constant.AuthUserFromCtx(ctx)
	if authUser == nil {
		return nil, sharedErrs.UnauthorizedErr
	}

	postExists, err := b.repository.FindPostByID(ctx, input.PostID, false)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}
	if postExists == nil {
		return nil, sharedErrs.NewBusinessValidationErr("Post %d not found.", input.PostID)
	}

	now := time.Now()
	comment, err := b.repository.CreateComment(ctx, &domain.Comment{
		BaseModel: domain.BaseModel{
			CreatedAt: now,
			UpdatedAt: now,
		},
		UserID:  authUser.ID,
		PostID:  input.PostID,
		Comment: input.Comment,
		Rating:  input.Rating,
	})
	if err != nil {
		return nil, err
	}

	return response.NewCommentFromDomain(comment), nil
}
