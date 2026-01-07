package middleware

import (
	"base_project/util/constant"
	"base_project/util/logger"
	"context"
	"net/http"

	"github.com/google/uuid"
)

// CorrelationID - Middleware to add requestID to context
// if "X-Request-ID" is empty then set in context
// if "X-Request-ID" is empty then add this in header and context
func CorrelationID(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requestID := r.Header.Get(constant.XCorrelationIDKey)
		ctx := r.Context()

		if requestID == "" {
			logger.Debug(ctx, "Creating new request id.")

			// Generate uuid
			requestID = uuid.New().String()

			ctx = context.WithValue(ctx, constant.XCorrelationIDKey, requestID)

			newRequest := r.WithContext(ctx)
			newRequest.Header.Set(constant.XCorrelationIDKey, requestID)

			logger.Debug(ctx, "CorrelationID: %v", requestID)
			next.ServeHTTP(w, newRequest)
		} else {
			// Add requestID to context
			ctx = context.WithValue(ctx, constant.XCorrelationIDKey, requestID)

			logger.Info(ctx, "CorrelationID: %v", requestID)
			next.ServeHTTP(w, r.WithContext(ctx))
		}
	})
}
