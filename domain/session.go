package domain

import (
	"time"

	"github.com/google/uuid"
)

type Session struct {
	ID                   uint64
	SessionID            uuid.UUID
	UserID               uint64
	AccessToken          string
	AccessTokenExpiresAt time.Time
	AccessTokenCreatedAt time.Time
	UserAgent            string
	ClientIP             string
}
