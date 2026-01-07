package response

import (
	"github.com/google/uuid"
)

type User struct {
	ID       uint64 `json:"id"`
	Username string `json:"username"`
	FullName string `json:"full_name"`
	Role     string `json:"role"`
}

type Auth struct {
	User                 User      `json:"user"`
	SessionID            uuid.UUID `json:"session_id"`
	AccessToken          string    `json:"access_token"`
	AccessTokenExpiresAt int64     `json:"access_token_expires_at"`
}
