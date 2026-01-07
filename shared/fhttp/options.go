package fhttp

import (
	"errors"
	"net/http"
)

const (
	healthCheckURL                  = "/healthcheck"
	defaultStopTimeoutSeconds       = 2
	defaultReadTimeoutSeconds       = 10
	defaultReadHeaderTimeoutSeconds = 5
)

type serverOptions struct {
	methodNotAllowedHandler  http.Handler
	notFoundHandler          http.Handler
	serverStopTimeoutSeconds int
	readTimeoutSeconds       int
	readHeaderTimeoutSeconds int
}

type Option func(*serverOptions) error

// WithMethodNotAllowedHandler overwrite default method not found handler
func WithMethodNotAllowedHandler(handler http.Handler) Option {
	return func(option *serverOptions) error {
		if handler != nil {
			option.methodNotAllowedHandler = handler
			return nil
		}

		return errors.New("MethodNotAllowedHandler is nil")
	}
}

// WithNotFoundHandler Overwrite default url not found handler
func WithNotFoundHandler(handler http.Handler) Option {
	return func(option *serverOptions) error {
		if handler != nil {
			option.notFoundHandler = handler
			return nil
		}

		return errors.New("NotFoundHandler is nil")
	}
}

// WithTimeout : timeout to stop server gracefully
func WithTimeout(timeout int) Option {
	return func(option *serverOptions) error {
		if timeout > 0 {
			option.serverStopTimeoutSeconds = timeout
			return nil
		}

		return errors.New("timeout should be greater than zero")
	}
}

// WithReadTimeout adds read timeout to the server.
func WithReadTimeout(timeoutSeconds int) Option {
	return func(option *serverOptions) error {
		if timeoutSeconds > 0 {
			option.readTimeoutSeconds = timeoutSeconds
			return nil
		}

		return errors.New("read timeout should be greater than zero")
	}
}

// WithReadHeaderTimeout adds read header timeout to the server.
func WithReadHeaderTimeout(timeoutSeconds int) Option {
	return func(option *serverOptions) error {
		if timeoutSeconds > 0 {
			option.readHeaderTimeoutSeconds = timeoutSeconds
			return nil
		}

		return errors.New("read header timeout should be greater than zero")
	}
}
