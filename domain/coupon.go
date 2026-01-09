package domain

import (
	"gorm.io/gorm"
)

type Coupon struct {
	BaseModel

	gorm.DeletedAt
	Name            string
	Amount          uint64
	RemainingAmount uint64

	// Association
	ClaimBy []*User `gorm:"many2many:user_claims;foreignKey:ID;joinForeignKey:CouponID;References:ID;joinReferences:UserID"`
}

func (c *Coupon) IsUsable() bool {
	if c == nil {
		return false
	}

	if c.RemainingAmount <= 0 {
		return false
	}

	return true
}
