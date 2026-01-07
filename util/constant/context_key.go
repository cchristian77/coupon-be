package constant

import (
	"context"
)

type contextKey string

func (c contextKey) String() string {
	return string(c)
}

const (
	correlationIDKey = contextKey("Correlation-ID")
	idempotencyKey   = contextKey("Idempotency-Key")
	ipAddressKey     = contextKey("ip-address")
)

var (
	XCorrelationIDKey = correlationIDKey.String()
	XIPAddressKey     = ipAddressKey.String()
)

func CorrelationIDFromCtx(ctx context.Context) string {
	correlationID, ok := ctx.Value(XCorrelationIDKey).(string)
	if !ok {
		return ""
	}

	return correlationID
}
