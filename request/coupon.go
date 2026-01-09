package request

type FilterCoupon struct {
	Search string `json:"search"`
}

type UpsertCoupon struct {
	ID uint64 `json:"-"`

	Name   string `json:"coupon_name" validate:"required"`
	Amount uint64 `json:"amount" validate:"required,gt=0"`
}

type ClaimCoupon struct {
	UserName   string `json:"user_id" validate:"required"`
	CouponName string `json:"coupon_name" validate:"required"`
}
