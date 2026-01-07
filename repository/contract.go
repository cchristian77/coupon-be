package repository

import (
	"base_project/domain"
	"base_project/domain/enums"
	"base_project/util"
	"context"
)

//go:generate mockgen -package mock -source=contract.go -destination=../mock/repository.go *

type Repository interface {
	// Session
	CreateSession(ctx context.Context, data *domain.Session) (*domain.Session, error)
	FindSessionBySessionID(ctx context.Context, sessionID string) (*domain.Session, error)
	DeleteSessionByID(ctx context.Context, id uint64) error

	// User
	FindUserByUsername(ctx context.Context, username string) (*domain.User, error)
	FindUserByEmail(ctx context.Context, email string) (*domain.User, error)
	FindUserByID(ctx context.Context, id uint64) (*domain.User, error)
	FindBatchUsers(ctx context.Context, batchSize int, lastID uint64) ([]*domain.User, error)
	CreateUser(ctx context.Context, data *domain.User) (*domain.User, error)

	// Post
	FindPublishedPostsPaginated(ctx context.Context, search string, p *util.Pagination) ([]*domain.Post, error)
	FindPostsPaginatedByUserID(ctx context.Context, userID uint64, search string, p *util.Pagination) ([]*domain.Post, error)
	FindPostByID(ctx context.Context, id uint64, withComments bool) (*domain.Post, error)
	FindPostByIDAndUserID(ctx context.Context, id, userID uint64) (*domain.Post, error)
	FindPostBySlug(ctx context.Context, slug string) (*domain.Post, error)
	CreatePost(ctx context.Context, data *domain.Post) (*domain.Post, error)
	UpdatePost(ctx context.Context, data *domain.Post) (*domain.Post, error)
	UpdatePostStatus(ctx context.Context, id uint64, status enums.PostStatus) (*domain.Post, error)
	DeletePostByID(ctx context.Context, id uint64) error

	// Comment
	FindCommentsPaginatedByPostID(ctx context.Context, postID uint64, p *util.Pagination) ([]*domain.Comment, error)
	FindCommentsPaginatedByUserID(ctx context.Context, userID uint64, p *util.Pagination) ([]*domain.Comment, error)
	FindCommentByID(ctx context.Context, id uint64) (*domain.Comment, error)
	FindCommentByIDAndUserID(ctx context.Context, id, userID uint64) (*domain.Comment, error)
	CreateComment(ctx context.Context, data *domain.Comment) (*domain.Comment, error)
	UpdateComment(ctx context.Context, data *domain.Comment) (*domain.Comment, error)
	DeleteCommentByID(ctx context.Context, id uint64) error
}
