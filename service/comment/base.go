package comment

import (
	"context"
	"coupon_be/repository"
	"coupon_be/request"
	"coupon_be/response"

	"gorm.io/gorm"
)

//go:generate mockgen -package service -source=base.go -destination=../../mock/service/comment.go -mock_names Service=MockCommentService *

type Service interface {
	FilterComments(ctx context.Context, input *request.FilterComment) (*response.BasePagination[[]*response.Comment], error)
	Detail(ctx context.Context, id uint64) (*response.Comment, error)
	Store(ctx context.Context, input *request.CreateComment) (*response.Comment, error)
	Update(ctx context.Context, input *request.UpdateComment) (*response.Comment, error)
	Delete(ctx context.Context, id uint64) error
}

type base struct {
	repository repository.Repository
	writeDB    *gorm.DB
}

func NewService(repository repository.Repository, writerDB *gorm.DB) (Service, error) {
	return &base{
		repository: repository,
		writeDB:    writerDB,
	}, nil
}
