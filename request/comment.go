package request

type CreateComment struct {
	PostID  uint64 `json:"post_id" validate:"required"`
	Comment string `json:"comment" validate:"required"`
	Rating  *uint8 `json:"rating"`
}

type UpdateComment struct {
	ID uint64 `json:"-"`

	Comment string `json:"comment" validate:"required"`
	Rating  *uint8 `json:"rating"`
}

type FilterComment struct {
	PostID uint64 `json:"-"`

	Page    int    `json:"page"`
	PerPage int    `json:"per_page"`
	Search  string `json:"search"`
}
