package auth

import (
	"base_project/repository"
	"base_project/service/auth"
	"base_project/util/logger"
	"context"
	"fmt"

	"gorm.io/gorm"
)

// NewController initializes a new Controller instance.
func NewController(ctx context.Context, repository repository.Repository, writerDB *gorm.DB) (*Controller, error) {
	authService, err := auth.NewService(repository, writerDB)
	if err != nil {
		logger.L().Fatal(fmt.Sprintf("auth service initialization error: %v", err))
	}

	return &Controller{auth: authService}, nil
}
