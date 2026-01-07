package comment

import (
	"context"
	"coupon_be/repository"
	"coupon_be/service/comment"
	"coupon_be/util/logger"
	"fmt"

	"gorm.io/gorm"
)

// NewController initializes a new Controller instance.
func NewController(ctx context.Context, repository repository.Repository, writerDB *gorm.DB) (*Controller, error) {
	commentService, err := comment.NewService(repository, writerDB)
	if err != nil {
		logger.L().Fatal(fmt.Sprintf("comment service initialization error: %v", err))
	}

	return &Controller{comment: commentService}, nil
}
