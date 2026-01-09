package response

import "coupon_be/domain"

type Coupon struct {
	Name            string   `json:"name"`
	Amount          uint64   `json:"amount"`
	RemainingAmount uint64   `json:"remaining_amount"`
	ClaimBy         []string `json:"claim_by"`
}

type CouponList struct {
	Name   string `json:"name"`
	Amount uint64 `json:"amount"`
}

func NewCouponFromDomain(c *domain.Coupon) *Coupon {
	if c == nil {
		return nil
	}

	claimBy := make([]string, len(c.ClaimBy))
	for i, claim := range c.ClaimBy {
		claimBy[i] = claim.Username
	}

	return &Coupon{
		Name:            c.Name,
		Amount:          c.Amount,
		RemainingAmount: c.RemainingAmount,
		ClaimBy:         claimBy,
	}
}

func NewCouponListFromDomains(c []*domain.Coupon) []*CouponList {
	result := make([]*CouponList, len(c))
	for i, coupon := range c {
		result[i] = &CouponList{
			Name:   coupon.Name,
			Amount: coupon.Amount,
		}
	}

	return result
}
