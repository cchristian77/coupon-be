package auth

import (
	"base_project/domain"
	"base_project/repository"
	"base_project/request"
	"base_project/response"
	"base_project/util/token"
	"context"

	"gorm.io/gorm"
)

type Service interface {
	Authenticate(ctx context.Context, accessToken string) (*domain.User, *token.Payload, error)
	Login(ctx context.Context, input *request.Login) (*response.Auth, error)
	Logout(ctx context.Context) error
	Register(ctx context.Context, input *request.Register) (*response.User, error)
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
