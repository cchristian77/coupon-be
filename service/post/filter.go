package post

import (
	"base_project/request"
	"base_project/response"
	"base_project/util"
	"context"
)

func (b *base) FilterPosts(ctx context.Context, input *request.FilterPost) (*response.BasePagination[[]*response.Post], error) {
	var p util.Pagination
	p.SetPage(input.Page)
	p.SetLimit(input.PerPage)

	posts, err := b.repository.FindPublishedPostsPaginated(ctx, input.Search, &p)
	if err != nil {
		return nil, err
	}

	result := make([]*response.Post, len(posts))
	for i, post := range posts {
		result[i] = response.NewPostFromDomain(post)
	}

	return response.NewBasePagination(result, &p), nil
}
