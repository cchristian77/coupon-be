package response

import (
	"coupon_be/domain"
	"coupon_be/util"
	"time"
)

type Comment struct {
	ID        uint64    `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`

	Username string `json:"username"`
	Comment  string `json:"comment"`
	Rating   *uint8 `json:"rating"`
}

func NewCommentFromDomain(c *domain.Comment) *Comment {
	if c == nil {
		return nil
	}

	return &Comment{
		ID:        c.ID,
		CreatedAt: c.CreatedAt,
		UpdatedAt: c.UpdatedAt,
		Username:  util.GetPointerValue(c.User).Username,
		Comment:   c.Comment,
		Rating:    c.Rating,
	}
}
