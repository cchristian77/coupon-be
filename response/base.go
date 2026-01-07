package response

import "base_project/util"

// Meta represents metadata for paginated responses.
type Meta struct {
	Page      int   `json:"page,omitempty"`
	PerPage   int   `json:"per_page,omitempty"`
	PageCount int   `json:"page_count"`
	Total     int64 `json:"total"`
}

// BasePagination represents a generic paginated response structure containing data and pagination metadata.
type BasePagination[T any] struct {
	Data T     `json:"data"`
	Meta *Meta `json:"meta"`
}

func NewBasePagination[T any](data []T, p *util.Pagination) *BasePagination[[]T] {
	return &BasePagination[[]T]{
		Data: data,
		Meta: &Meta{
			Page:      p.Page(),
			PerPage:   len(data),
			PageCount: p.PageCount(),
			Total:     p.Total(),
		},
	}

}
