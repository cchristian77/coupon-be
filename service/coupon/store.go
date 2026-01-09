package coupon

import (
	"context"
	"coupon_be/domain"
	"coupon_be/request"
	"coupon_be/response"
	sharedErrs "coupon_be/shared/errors"
	"coupon_be/util/logger"
	"errors"
	"regexp"
	"strings"
	"time"
)

func (b base) Store(ctx context.Context, input *request.UpsertCoupon) (*response.Coupon, error) {
	logger.Info(ctx, "Store Post with request: %v", input)

	input.Name = toCouponName(input.Name)

	couponExists, err := b.repository.FindCouponByName(ctx, input.Name, false)
	if err != nil && !errors.Is(err, sharedErrs.NotFoundErr) {
		return nil, err
	}
	if couponExists != nil {
		return nil, sharedErrs.NewBusinessValidationErr(
			"Create Failed. Coupon with name '%s' already exists.", input.Name)
	}

	now := time.Now()
	coupon, err := b.repository.CreateCoupon(ctx, &domain.Coupon{
		BaseModel: domain.BaseModel{
			CreatedAt: now,
			UpdatedAt: now,
		},
		Name:            input.Name,
		Amount:          input.Amount,
		RemainingAmount: input.Amount,
	})
	if err != nil {
		return nil, err
	}

	return response.NewCouponFromDomain(coupon), nil
}

func toCouponName(s string) string {
	name := strings.ToUpper(s)

	// Replace any character that is NOT a-z or 0-9 with a underscore
	// The regex pattern `[^A-Z0-9]+` matches one or more non-alphanumeric characters
	reg := regexp.MustCompile(`[^A-Z0-9]+`)
	name = reg.ReplaceAllString(name, "_")

	// rim hyphens from the start and end
	name = strings.Trim(name, "_")

	return name
}
