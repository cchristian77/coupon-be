package response

import (
	"base_project/domain"
	"base_project/domain/enums"
	"time"
)

type Post struct {
	ID        uint64    `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`

	Slug   string           `json:"slug"`
	Title  string           `json:"title"`
	Body   string           `json:"body"`
	Status enums.PostStatus `json:"status"`

	*Author `json:"author,omitempty"`
}

type Author struct {
	Username string `json:"username"`
	FullName string `json:"full_name"`
}

func NewPostFromDomain(p *domain.Post) *Post {
	if p == nil {
		return nil
	}

	var author *Author
	if p.User != nil {
		author = &Author{
			Username: p.User.Username,
			FullName: p.User.FullName,
		}
	}

	return &Post{
		ID:        p.ID,
		CreatedAt: p.CreatedAt,
		UpdatedAt: p.UpdatedAt,
		Slug:      p.Slug,
		Title:     p.Title,
		Body:      p.Body,
		Status:    p.Status,
		Author:    author,
	}
}
