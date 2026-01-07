package domain

import (
	"base_project/domain/enums"
	"time"
)

type User struct {
	BaseModel

	DeletedAt *time.Time
	Username  string
	Email     string
	FullName  string
	Password  string
	Role      enums.UserRole
}
