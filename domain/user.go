package domain

import "gorm.io/gorm"

type User struct {
	BaseModel

	DeletedAt gorm.DeletedAt
	Username  string
	FullName  string
	Password  string
}
