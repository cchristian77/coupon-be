package middleware

import (
	"base_project/domain/enums"
	"base_project/repository"
	"base_project/service/auth"
	sharedErrs "base_project/shared/errors"
	"base_project/shared/fhttp"
	"base_project/util"
	"base_project/util/constant"
	"base_project/util/logger"
	"context"
	"fmt"
	"net/http"
	"strings"

	"gorm.io/gorm"
)

var authMiddleware *Authorization

// authHeaderKey represents the key for the Authorization header in HTTP requests.
// authTypeBearer defines the bearer authentication type used in the Authorization header.
const (
	authHeaderKey  = "Authorization"
	authTypeBearer = "bearer"
)

// Authorization is a middleware model that provides role-based authorization for HTTP requests.
type Authorization struct {
	authService auth.Service
}

// NewAuthMiddleware initializes the authorization middleware with the provided authentication service.
func NewAuthMiddleware(ctx context.Context, repo repository.Repository, writerDB *gorm.DB) {
	authService, err := auth.NewService(repo, writerDB)
	if err != nil {
		logger.L().Fatal(fmt.Sprintf("auth service initialization error: %v", err))
	}

	authMiddleware = &Authorization{
		authService: authService,
	}

	return
}

// AdminOnly is an HTTP middleware that restricts access only to users with the ADMIN role.
func AdminOnly() func(h http.Handler) http.Handler {
	return authMiddleware.authWithRoles(enums.ADMINRole.String())
}

// UserOnly is a middleware function that restricts access only to users with the USER role.
func UserOnly() func(h http.Handler) http.Handler {
	return authMiddleware.authWithRoles(enums.USERRole.String())
}

// Authentication provides middleware for authorizing requests by validating the user's role as either ADMIN or USER.
func Authentication() func(h http.Handler) http.Handler {
	return authMiddleware.authWithRoles(enums.ADMINRole.String(), enums.USERRole.String())
}

// authenticationWithRoles creates middleware to authenticate users and check if they have one of the specified roles.
func (a *Authorization) authWithRoles(allowedRoles ...string) func(h http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := r.Context()

			authHeader := r.Header.Get(authHeaderKey)
			if authHeader == "" {
				fhttp.WriteErrorResponse(ctx, sharedErrs.UnauthorizedErr, w)
				return
			}

			authFields := strings.Fields(authHeader)
			if len(authFields) < 2 {
				fhttp.WriteErrorResponse(ctx, sharedErrs.UnauthorizedErr, w)
				return
			}

			authorizationType := strings.ToLower(authFields[0])
			if authorizationType != authTypeBearer {
				fhttp.WriteErrorResponse(ctx, sharedErrs.UnauthorizedErr, w)
				return
			}

			bearerToken := authFields[1]

			// authenticate the bearer token
			authUser, payload, err := a.authService.Authenticate(ctx, bearerToken)
			if err != nil {
				fhttp.WriteErrorResponse(ctx, err, w)
				return
			}

			// authorization based on the allowed roles from the authenticated user's role
			// find if user's role contains in the allowed roles
			if util.Contains(allowedRoles, authUser.Role.String()) {
				ctx = context.WithValue(ctx, constant.XAuthUserKey, authUser)
				ctx = context.WithValue(ctx, constant.XSessionIDKey, payload.ID.String())
				ctx = context.WithValue(ctx, constant.XIPAddressKey, r.RemoteAddr)
				newRequest := r.WithContext(ctx)

				next.ServeHTTP(w, newRequest)
				return
			}

			// if not found, then return forbidden access error
			fhttp.WriteErrorResponse(ctx, sharedErrs.ForbiddenErr, w)
		})
	}
}
