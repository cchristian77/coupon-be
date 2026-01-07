package constant

import (
	"base_project/domain"
	"context"
)

type contextKey string

func (c contextKey) String() string {
	return string(c)
}

const (
	correlationIDKey = contextKey("Correlation-ID")
	eventIDKey       = contextKey("Event-ID")
	idempotencyKey   = contextKey("Idempotency-Key")
	authUserKey      = contextKey("auth-user")
	userIDKey        = contextKey("user-id")
	sessionIDKey     = contextKey("session-id")
	ipAddressKey     = contextKey("ip-address")
)

var (
	XCorrelationIDKey = correlationIDKey.String()
	XAuthUserKey      = authUserKey.String()
	XUserIDKey        = userIDKey.String()
	XSessionIDKey     = sessionIDKey.String()
	XIPAddressKey     = ipAddressKey.String()
)

func AuthUserFromCtx(ctx context.Context) *domain.User {
	authUser, ok := ctx.Value(XAuthUserKey).(*domain.User)
	if !ok {
		return nil
	}

	return authUser
}

func CorrelationIDFromCtx(ctx context.Context) string {
	correlationID, ok := ctx.Value(XCorrelationIDKey).(string)
	if !ok {
		return ""
	}

	return correlationID
}
