package post

import (
	"base_project/domain/enums"
	"base_project/repository"
	"base_project/request"
	"base_project/response"
	"context"

	"gorm.io/gorm"
)

//go:generate mockgen -package service -source=base.go -destination=../../mock/service/post.go -mock_names Service=MockPostService *

type Service interface {
	FilterPosts(ctx context.Context, input *request.FilterPost) (*response.BasePagination[[]*response.Post], error)
	FilterMyPosts(ctx context.Context, input *request.FilterPost) (*response.BasePagination[[]*response.Post], error)
	Detail(ctx context.Context, id uint64) (*response.Post, error)
	Store(ctx context.Context, input *request.UpsertPost) (*response.Post, error)
	Update(ctx context.Context, input *request.UpsertPost) (*response.Post, error)
	UpdateStatus(ctx context.Context, id uint64, status enums.PostStatus) error
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
