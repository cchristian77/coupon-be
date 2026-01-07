package fhttp

import (
	"fmt"
	"net/http"
)

// DefaultHealthCheckHandler Default handler for health check
func DefaultHealthCheckHandler(r *http.Request) (*Response, error) {
	return &Response{
		Data:   "Service is Running.",
		Status: http.StatusOK,
	}, nil
}

// DefaultMethodNotAllowedHandler default handler for method not allowed
func DefaultMethodNotAllowedHandler(r *http.Request) (*Response, error) {
	return &Response{
		Data:   fmt.Sprintf("%v Method not allowed", r.Method),
		Status: http.StatusMethodNotAllowed,
	}, nil
}

// DefaultNotFoundHandler default handler for not found handler
func DefaultNotFoundHandler(r *http.Request) (*Response, error) {
	return &Response{
		Data:   "Url not found",
		Status: http.StatusNotFound,
	}, nil
}
