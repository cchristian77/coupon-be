package coupon

import (
	"context"
	"coupon_be/request"
	"coupon_be/response"
	"coupon_be/util/logger"
)

func (b base) Filter(ctx context.Context, input *request.FilterCoupon) ([]*response.CouponList, error) {
	logger.Info(ctx, "Filter Coupon with req: %d", input)

	coupons, err := b.repository.FindCoupons(ctx, input.Search)
	if err != nil {
		return nil, err
	}

	return response.NewCouponListFromDomains(coupons), nil
}
