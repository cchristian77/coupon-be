package middleware

import (
	"base_project/shared/fhttp"
	"base_project/util/logger"
	"fmt"
	"net/http"
	"runtime/debug"
)

func PanicRecovery() func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			defer func() {
				ctx := r.Context()
				if p := recover(); p != nil {
					err := fmt.Errorf("Panic: %v, Panic stacktrace: \n %v", p, string(debug.Stack()))
					logger.Error(ctx, "%s", err.Error())

					fhttp.WriteErrorResponse(ctx, err, w)
				}
			}()

			next.ServeHTTP(w, r)
		})
	}
}
