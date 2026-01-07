package user

import (
	"context"
	"coupon_be/repository"

	"gorm.io/gorm"
)

type Service interface {
	Register(ctx context.Context) error
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
