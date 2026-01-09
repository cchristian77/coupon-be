package errors

// Kind is used to pinpoint a specific error.
type Kind string

func (k Kind) String() string {
	return string(k)
}

// Application Error Kinds
const (
	// HTTP or Controller Errors
	ErrKindHTTP           Kind = "http_error"
	ErrKindInvalidRequest Kind = "invalid_request_error"
	ErrKindValidation     Kind = "validation_error"

	// ErrKindApplication errors are errors that might be resolved by retrying the same request at a later time
	ErrKindApplication Kind = "application_error"
	// ErrKindApplicationPermanent Permanent errors are errors that are not expected to be resolved by retrying the same request
	ErrKindApplicationPermanent Kind = "application_permanent_error"
	ErrKindBusinessValidation   Kind = "business_validation_error"
	ErrKindCodeInjection        Kind = "code_injection_error"

	// Repository Error
	ErrKindRepository   Kind = "repository_error"
	ErrKindDatabase     Kind = "database_error"
	ErrKindDataNotFound Kind = "data_not_found_error"

	// Utilities or Misclanneous Errors
	ErrKindAuthorization Kind = "authorization_error"
	ErrKindForbidden     Kind = "forbidden_error"
	ErrKindConflict      Kind = "conflict_error"
	ErrKindUnknown       Kind = "unknown_error"
)

// External Client Error Kinds
const (
	ErrKindClientExternal   Kind = "external_client_error"
	ErrKindHttpClient       Kind = "http_client_error"
	ErrKindRedis            Kind = "redis_error"
	ErrKindAcquireRedisLock Kind = "acquire_redis_lock_error"
	ErrKindDependency       Kind = "dependency_error"
)
