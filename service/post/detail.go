package post

import (
	"base_project/response"
	sharedErrs "base_project/shared/errors"
	"base_project/util/constant"
	"context"
)

func (b *base) Detail(ctx context.Context, id uint64) (*response.Post, error) {
	authUser := constant.AuthUserFromCtx(ctx)
	if authUser == nil {
		return nil, sharedErrs.UnauthorizedErr
	}

	post, err := b.repository.FindPostByID(ctx, id, false)
	if err != nil {
		return nil, err
	}

	if !post.IsPublished() && authUser.ID != post.UserID {
		return nil, sharedErrs.NotFoundErr
	}

	return response.NewPostFromDomain(post), nil
}
