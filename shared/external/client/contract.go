package client

import "context"

type Client interface {
	CallAPI(ctx context.Context, method, url string, reqBody interface{},
		headers map[string]string,
	) (rawRespBody string, httpStatusCode int, err error)
}
