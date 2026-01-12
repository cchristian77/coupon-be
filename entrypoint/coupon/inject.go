package coupon

import (
	"context"
	"coupon_be/repository"
	"coupon_be/service/coupon"
	sharedErrs "coupon_be/shared/errors"
	"coupon_be/shared/external/redis"

	"gorm.io/gorm"
)

// NewController initializes a new Controller instance.
func NewController(ctx context.Context, repository repository.Repository, writerDB *gorm.DB) (*Controller, error) {
	couponService, err := coupon.NewService(repository, writerDB)
	if err != nil {
		return nil, sharedErrs.NewWithCause(sharedErrs.ErrKindCodeInjection, "Fail to initiate coupon service", err)
	}

	redisLock, err := redis.GetRedisLock(ctx)
	if err != nil {
		return nil, sharedErrs.NewWithCause(sharedErrs.ErrKindApplication, "Fail to initiate Redis Lock", err)
	}

	return &Controller{coupon: couponService, redisLock: redisLock}, nil
}
