package coupon

import (
	"context"
	"coupon_be/request"
	"coupon_be/response"
	"coupon_be/util"
	"coupon_be/util/logger"
)

func (b *base) Filter(ctx context.Context, input *request.FilterCoupon) (*response.BasePagination[[]*response.CouponList], error) {
	logger.Info(ctx, "Filter Coupon with req: %v", input)

	var p util.Pagination
	p.SetPage(input.Page)
	p.SetLimit(input.PerPage)

	coupons, err := b.repository.FindCouponsPaginated(ctx, input.Search, &p)
	if err != nil {
		return nil, err
	}

	result := make([]*response.CouponList, len(coupons))
	for i, coupon := range coupons {
		result[i] = response.NewCouponListFromDomain(coupon)
	}

	return response.NewBasePagination(result, &p), nil
}
