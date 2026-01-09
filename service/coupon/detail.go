package coupon

import (
	"context"
	"coupon_be/response"
	"coupon_be/util"
	"coupon_be/util/logger"
)

func (b base) Detail(ctx context.Context, name string) (*response.Coupon, error) {
	logger.Info(ctx, "Get Detail Coupon with name: %s", name)

	coupon, err := b.repository.FindCouponByName(ctx, util.SanitizeString(name), true)
	if err != nil {
		return nil, err
	}

	return response.NewCouponFromDomain(coupon), nil
}
