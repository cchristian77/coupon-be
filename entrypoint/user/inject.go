package user

import (
	"context"
	"coupon_be/repository"
	"coupon_be/service/user"
	sharedErrs "coupon_be/shared/errors"

	"gorm.io/gorm"
)

// NewController initializes a new Controller instance.
func NewController(ctx context.Context, repository repository.Repository, writerDB *gorm.DB) (*Controller, error) {
	authService, err := user.NewService(repository, writerDB)
	if err != nil {
		return nil, sharedErrs.NewWithCause(sharedErrs.ErrKindCodeInjection, "Fail to initiate user service", err)
	}

	return &Controller{user: authService}, nil
}
