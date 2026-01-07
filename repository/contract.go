package repository

import (
	"context"
	"coupon_be/domain"
	"coupon_be/util"
)

//go:generate mockgen -package mock -source=contract.go -destination=../mock/repository.go *

type Repository interface {
	// User
	FindUserByUsername(ctx context.Context, username string) (*domain.User, error)
	FindUserByID(ctx context.Context, id uint64) (*domain.User, error)
	CreateUser(ctx context.Context, data *domain.User) (*domain.User, error)

	// Comment
	FindCommentsPaginatedByPostID(ctx context.Context, postID uint64, p *util.Pagination) ([]*domain.Comment, error)
	FindCommentsPaginatedByUserID(ctx context.Context, userID uint64, p *util.Pagination) ([]*domain.Comment, error)
	FindCommentByID(ctx context.Context, id uint64) (*domain.Comment, error)
	FindCommentByIDAndUserID(ctx context.Context, id, userID uint64) (*domain.Comment, error)
	CreateComment(ctx context.Context, data *domain.Comment) (*domain.Comment, error)
	UpdateComment(ctx context.Context, data *domain.Comment) (*domain.Comment, error)
	DeleteCommentByID(ctx context.Context, id uint64) error
}
