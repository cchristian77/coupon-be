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

	// Coupon
	FindCouponByID(ctx context.Context, id uint64) (*domain.Coupon, error)
	FindCouponByName(ctx context.Context, name string, withClaimBy bool) (*domain.Coupon, error)
	FindCouponsPaginated(ctx context.Context, search string, p *util.Pagination) ([]*domain.Coupon, error)
	CreateCoupon(ctx context.Context, data *domain.Coupon) (*domain.Coupon, error)
	UpdateCoupon(ctx context.Context, data *domain.Coupon) (*domain.Coupon, error)
	DecrementCouponRemainingAmount(ctx context.Context, id uint64) error

	// User Claim
	FindUserClaimByUserIDAndCouponID(ctx context.Context, userID, couponID uint64) (*domain.UserClaim, error)
	FindUserClaimCountByCouponID(ctx context.Context, couponID uint64) (int64, error)
	CreateUserClaim(ctx context.Context, data *domain.UserClaim) (*domain.UserClaim, error)
}
