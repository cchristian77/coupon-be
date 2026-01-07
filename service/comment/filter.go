package comment

import (
	"base_project/request"
	"base_project/response"
	sharedErrs "base_project/shared/errors"
	"base_project/util"
	"context"
	"errors"

	"gorm.io/gorm"
)

func (b *base) FilterComments(ctx context.Context, input *request.FilterComment) (*response.BasePagination[[]*response.Comment], error) {
	var p util.Pagination
	p.SetPage(input.Page)
	p.SetLimit(input.PerPage)

	postExists, err := b.repository.FindPostByID(ctx, input.PostID, false)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}
	if postExists == nil {
		return nil, sharedErrs.NewBusinessValidationErr("Post %d not found.", input.PostID)
	}

	comments, err := b.repository.FindCommentsPaginatedByPostID(ctx, input.PostID, &p)
	if err != nil {
		return nil, err
	}

	result := make([]*response.Comment, len(comments))
	for i, comment := range comments {
		result[i] = response.NewCommentFromDomain(comment)
	}

	return response.NewBasePagination(result, &p), nil
}
