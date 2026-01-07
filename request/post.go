package request

type FilterPost struct {
	Page    int    `json:"page"`
	PerPage int    `json:"per_page"`
	Search  string `json:"search"`
}

type UpsertPost struct {
	ID uint64 `json:"-"`

	Slug  string `json:"slug" validate:"required"`
	Title string `json:"title" validate:"required,min=6"`
	Body  string `json:"body" validate:"required"`
}

type UpdatePostStatus struct {
	Status string `json:"status"`
}
