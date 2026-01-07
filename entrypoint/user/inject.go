package user

import (
	"context"
	"coupon_be/repository"
	"coupon_be/service/user"
	"coupon_be/util/logger"
	"fmt"

	"gorm.io/gorm"
)

// NewController initializes a new Controller instance.
func NewController(ctx context.Context, repository repository.Repository, writerDB *gorm.DB) (*Controller, error) {
	authService, err := user.NewService(repository, writerDB)
	if err != nil {
		logger.L().Fatal(fmt.Sprintf("auth service initialization error: %v", err))
	}

	return &Controller{auth: authService}, nil
}
