package comment

import (
	"context"
	"coupon_be/domain"
	"coupon_be/request"
	"coupon_be/response"
	sharedErrs "coupon_be/shared/errors"
	"coupon_be/util/constant"
	"coupon_be/util/logger"
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
