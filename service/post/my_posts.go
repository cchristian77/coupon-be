package post

import (
	"base_project/request"
	"base_project/response"
	sharedErrs "base_project/shared/errors"
	"base_project/util"
	"base_project/util/constant"
	"context"
)

func (b *base) FilterMyPosts(ctx context.Context, input *request.FilterPost) (*response.BasePagination[[]*response.Post], error) {
	authUser := constant.AuthUserFromCtx(ctx)
	if authUser == nil {
		return nil, sharedErrs.UnauthorizedErr
	}

	var p util.Pagination
	p.SetPage(input.Page)
	p.SetLimit(input.PerPage)

	posts, err := b.repository.FindPostsPaginatedByUserID(ctx, authUser.ID, input.Search, &p)
	if err != nil {
		return nil, err
	}

	result := make([]*response.Post, len(posts))
	for i, post := range posts {
		result[i] = response.NewPostFromDomain(post)
	}

	return response.NewBasePagination(result, &p), nil
}
