package token

import (
	sharedErrs "base_project/shared/errors"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/google/uuid"
)

const issuer = "payroll.api.auth"

// Payload contains the payload data of the token
type Payload struct {
	ID     uuid.UUID `json:"id"`
	UserID uint64    `json:"user_id"`
	jwt.StandardClaims
}

// NewPayload creates a new token payload with a specific username and duration
func NewPayload(sessionID uuid.UUID, userID uint64, duration time.Duration) (*Payload, error) {
	return &Payload{
		ID:     sessionID,
		UserID: userID,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(duration).Unix(),
			IssuedAt:  time.Now().Unix(),
			Issuer:    issuer,
		},
	}, nil
}

// Valid checks if the token payload is valid or not
func (payload *Payload) Valid() error {
	expiredAt := time.Unix(payload.StandardClaims.ExpiresAt, 0)

	if time.Now().After(expiredAt) {
		return sharedErrs.InvalidTokenErr
	}

	return nil
}
