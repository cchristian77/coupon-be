package errors

import (
	"errors"
	"net/http"

	"gorm.io/gorm"
)

// HTTP ERROR
var (
	// InternalServerErr will throw if any the Internal Server Error happen
	InternalServerErr = New(ErrKindHTTP, "Internal Server Error")

	// NotFoundErr will throw if the requested item is not exists
	NotFoundErr = New(ErrKindDataNotFound, "Requested data is not found")

	// ConflictErr will throw if the current action already exists
	ConflictErr = New(ErrKindRepository, "Requested data already exist")

	// BadParamInputErr will throw if the given request-body or params is not valid
	BadParamInputErr = New(ErrKindBusinessValidation, "Requested parameters are not valid")

	// ForbiddenErr will throw if the current request is forbidden
	ForbiddenErr = New(ErrKindForbidden, "Forbidden Access")
)

// AUTH ERROR
var (
	// UnauthorizedErr will throw if the current request is unauthorized
	UnauthorizedErr        = New(ErrKindAuthorization, "Unauthorized")
	IncorrectCredentialErr = New(ErrKindValidation, "Login failed. Email or password is incorrect.")
	InvalidTokenErr        = New(ErrKindAuthorization, "Invalid token")
	ExpiredTokenErr        = New(ErrKindAuthorization, "Expired token")
)

func GetStatusCode(err error) int {
	if err == nil {
		return http.StatusOK
	}

	if errors.Is(err, gorm.ErrRecordNotFound) {
		err = NotFoundErr
	}

	if baseErr, ok := err.(*baseError); ok {
		return getStatusCodeByErrorKind(baseErr.kind)
	}

	return http.StatusInternalServerError
}

func getStatusCodeByErrorKind(k Kind) int {
	switch k {
	case ErrKindValidation, ErrKindBusinessValidation, ErrKindApplication, ErrKindApplicationPermanent, ErrKindRepository:
		return http.StatusBadRequest
	case ErrKindDataNotFound:
		return http.StatusNotFound
	case ErrKindDatabase, ErrKindHTTP, ErrKindClientExternal, ErrKindRedis, ErrKindDependency:
		return http.StatusInternalServerError
	case ErrKindInvalidRequest:
		return http.StatusUnprocessableEntity
	case ErrKindAuthorization:
		return http.StatusUnauthorized
	case ErrKindForbidden:
		return http.StatusForbidden
	case ErrKindConflict, ErrKindAcquireRedisLock:
		return http.StatusConflict
	default:
		return http.StatusInternalServerError
	}
}
