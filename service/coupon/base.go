package coupon

import (
	"context"
	"coupon_be/repository"
	"coupon_be/request"
	"coupon_be/response"

	"gorm.io/gorm"
)

type Service interface {
	Filter(ctx context.Context, input *request.FilterCoupon) ([]*response.CouponList, error)

	Detail(ctx context.Context, name string) (*response.Coupon, error)

	Store(ctx context.Context, input *request.UpsertCoupon) (*response.Coupon, error)

	Claim(ctx context.Context, input *request.ClaimCoupon) error
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
