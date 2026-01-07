package domain

import (
	"base_project/domain/enums"

	"gorm.io/gorm"
)

type Post struct {
	BaseModel

	DeletedAt gorm.DeletedAt
	UserID    uint64
	Slug      string
	Title     string
	Body      string
	Status    enums.PostStatus

	Comments []*Comment
	User     *User
}

func (p *Post) IsPublished() bool {
	if p == nil {
		return false
	}

	if p.Status == enums.PUBLISHEDPostStatus {
		return true
	}

	return false
}
