package client

import (
	sharedErrs "base_project/shared/errors"
	"base_project/util/logger"
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"
)

// CallAPI calls the given API with the given method, request body, and headers.
func (c *client) CallAPI(ctx context.Context, method, url string, reqBody interface{},
	headers map[string]string,
) (rawRespBody string, httpStatusCode int, err error) {
	endpoint := c.host.String() + url

	// convert request struct to bytes
	reqBodyBytes := new(bytes.Buffer)
	if err := json.NewEncoder(reqBodyBytes).Encode(reqBody); err != nil {
		return "", 0, sharedErrs.NewWithCause(
			sharedErrs.ErrKindHttpClient, "Error encoding to bytes", err)
	}

	req, err := http.NewRequestWithContext(ctx, method, endpoint, reqBodyBytes)
	if err != nil {
		return "", 0, sharedErrs.NewWithCause(
			sharedErrs.ErrKindHttpClient, "Failed to create HTTP request", err)
	}

	// refer to https://stackoverflow.com/questions/28046100/golang-http-concurrent-requests-post-eof
	req.Close = true
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", c.authKey)

	for key, value := range headers {
		req.Header.Set(key, value)
	}

	logger.Debug(ctx, "Call API {%v}, method {%v}", url, method)

	response, err := c.httpClient.Do(req)
	if err != nil {
		return "", 0, sharedErrs.NewWithCause(
			sharedErrs.ErrKindHttpClient, "Failed to send HTTP request", err)
	}

	logger.Debug(ctx, "Call to API {%v}, method {%v} successful", url, method)

	defer func() {
		if err := response.Body.Close(); err != nil {
			logger.Error(ctx, "failed to close response body, err: %v", err)
		}
	}()

	rawRespBody, err = convertRawBodyToString(response.Body)
	if err != nil {
		return "", 0, sharedErrs.NewWithCause(
			sharedErrs.ErrKindHttpClient, "Failed to convert response body to string", err)
	}

	logger.Debug(ctx, "API call {%v} {%v} response: HTTP Status {%v}, Error {%v}, Response body: %v",
		method, url, response.StatusCode, err, rawRespBody)

	return rawRespBody, response.StatusCode, nil
}

func convertRawBodyToString(body io.ReadCloser) (string, error) {
	strRequestBody, err := io.ReadAll(body)
	if err != nil {
		return "", err
	}

	return string(strRequestBody), nil
}
