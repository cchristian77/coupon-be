package fhttp

import (
	"base_project/util/constant"
	"context"
	"net/http"
)

// AppHandler - Wrapper for controller functions
type AppHandler func(*http.Request) (*Response, error)

func (fn AppHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	resp, err := fn(r) // execute handler func

	if err != nil {
		WriteErrorResponse(ctx, err, w)

		return
	}

	WriteHttpResponse(ctx, resp, w)
}

func setHeaders(ctx context.Context, headers HTTPHeaders, w http.ResponseWriter) {
	if _, ok := headers[ContentTypeKey]; !ok {
		w.Header().Set(ContentTypeKey, string(ContentTypeJSON))
	}

	if val, ok := ctx.Value(constant.XCorrelationIDKey).(string); ok {
		w.Header().Set(constant.XCorrelationIDKey, val)
	}

	for key, val := range headers {
		w.Header().Set(key, val)
	}
}
