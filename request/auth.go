package request

type Login struct {
	Username string `json:"username" validate:"required,min=5"`
	Password string `json:"password" validate:"required"`

	ClientIP  string `json:"-"`
	UserAgent string `json:"-"`
}

type Register struct {
	Username string `json:"username" validate:"required,min=5"`
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
	FullName string `json:"full_name" validate:"required,min=5"`
}
