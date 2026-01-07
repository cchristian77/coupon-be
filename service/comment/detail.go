package comment

import (
	"context"
	"coupon_be/response"
)

func (b *base) Detail(ctx context.Context, id uint64) (*response.Comment, error) {
	comment, err := b.repository.FindCommentByID(ctx, id)
	if err != nil {
		return nil, err
	}

	return response.NewCommentFromDomain(comment), nil
}
