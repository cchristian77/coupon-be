package response

import "coupon_be/domain"

type Coupon struct {
	Name            string   `json:"name"`
	Amount          uint64   `json:"amount"`
	RemainingAmount uint64   `json:"remaining_amount"`
	ClaimedBy       []string `json:"claimed_by"`
}

type CouponList struct {
	Name   string `json:"name"`
	Amount uint64 `json:"amount"`
}

func NewCouponFromDomain(c *domain.Coupon) *Coupon {
	if c == nil {
		return nil
	}

	claimedBy := make([]string, len(c.ClaimedBy))
	for i, claim := range c.ClaimedBy {
		claimedBy[i] = claim.Username
	}

	return &Coupon{
		Name:            c.Name,
		Amount:          c.Amount,
		RemainingAmount: c.RemainingAmount,
		ClaimedBy:       claimedBy,
	}
}

func NewCouponListFromDomain(c *domain.Coupon) *CouponList {
	if c == nil {
		return nil
	}

	return &CouponList{
		Name:   c.Name,
		Amount: c.Amount,
	}
}
