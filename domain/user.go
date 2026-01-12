package domain

import "gorm.io/gorm"

type User struct {
	BaseModel

	DeletedAt gorm.DeletedAt `gorm:"index"`
	Username  string
	FullName  string
	Password  string
}
