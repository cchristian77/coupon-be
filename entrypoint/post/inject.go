package post

import (
	"base_project/repository"
	"base_project/service/post"
	"base_project/util/logger"
	"context"
	"fmt"

	"gorm.io/gorm"
)

// NewController initializes a new Controller instance.
func NewController(ctx context.Context, repository repository.Repository, writerDB *gorm.DB) (*Controller, error) {
	postService, err := post.NewService(repository, writerDB)
	if err != nil {
		logger.L().Fatal(fmt.Sprintf("post service initialization error: %v", err))
	}

	return &Controller{post: postService}, nil
}
