package domain

type UserClaim struct {
	BaseModel

	UserID   uint64
	CouponID uint64
}
