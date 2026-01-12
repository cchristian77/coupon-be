package domain

import (
	"gorm.io/gorm"
)

type Coupon struct {
	BaseModel

	DeletedAt       gorm.DeletedAt `gorm:"index"`
	Name            string
	Amount          uint64
	RemainingAmount uint64

	// Association
	ClaimedBy []*User `gorm:"many2many:user_claims;"`
}

func (c *Coupon) IsUsable() bool {
	if c == nil {
		return false
	}

	return c.RemainingAmount > 0
}
