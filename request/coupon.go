package request

type FilterCoupon struct {
	Page    int    `json:"page"`
	PerPage int    `json:"per_page"`
	Search  string `json:"search"`
}

type UpsertCoupon struct {
	ID uint64 `json:"-"`

	Name   string `json:"coupon_name" validate:"required"`
	Amount uint64 `json:"amount" validate:"required,gt=0"`
}

type ClaimCoupon struct {
	Username   string `json:"user_id" validate:"required"`
	CouponName string `json:"coupon_name" validate:"required"`
}
