package util

import "math"

const (
	DefaultPage  = 1
	DefaultLimit = 20
)

// Pagination - struct to hold required data and provide methods for pagination calculations.
type Pagination struct {
	limit int
	page  int
	total int64
}

// Limit returns the number of items per page.
func (p *Pagination) Limit() int {
	if p.limit <= 0 {
		p.limit = DefaultLimit
	}
	return p.limit
}

// Page - returns the current page number.
func (p *Pagination) Page() int {
	if p.page <= 0 {
		p.page = DefaultPage
	}
	return p.page
}

// Total returns the total number of items in the dataset.
func (p *Pagination) Total() int64 {
	return p.total
}

// Offset calculates and returns the offset for the database query.
func (p *Pagination) Offset() int {
	return (p.Page() - 1) * p.Limit()
}

// PageCount returns the maximum number of pages available.
func (p *Pagination) PageCount() int {
	return int(math.Ceil(float64(p.total) / float64(p.limit)))
}

// SetLimit set limit attribute on pagination struct
func (p *Pagination) SetLimit(limit int) {
	p.limit = limit
}

// SetPage set page attribute on pagination struct
func (p *Pagination) SetPage(page int) {
	p.page = page
}

// SetTotal set total attribute on pagination struct
func (p *Pagination) SetTotal(total int64) {
	p.total = total
}
