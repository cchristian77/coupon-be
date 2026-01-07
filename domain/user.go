package domain

import "time"

type User struct {
	BaseModel

	DeletedAt *time.Time
	Username  string
	FullName  string
	Password  string
}
