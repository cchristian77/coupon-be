package client

import (
	sharedErrs "base_project/shared/errors"
	"net/http"
	"net/url"
	"time"
)

type client struct {
	httpClient *http.Client
	host       *url.URL
	authKey    string
	options    clientOptions
}

func NewClient(timeout time.Duration, host string, authKey string, opts ...ClientOption) (Client, error) {
	u, err := url.Parse(host)
	if err != nil {
		return nil, sharedErrs.NewWithCause(sharedErrs.ErrKindHttpClient, "failed to parse host", err)
	}

	c := &client{
		httpClient: &http.Client{Timeout: timeout},
		host:       u,
		authKey:    authKey,
	}

	for _, opt := range opts {
		if err = opt(&c.options); err != nil {
			return nil, sharedErrs.NewWithCause(sharedErrs.ErrKindHttpClient, "failed to apply client option", err)
		}
	}

	if c.options.httpClient != nil {
		c.httpClient = c.options.httpClient
	}

	if c.options.transport != nil {
		c.httpClient.Transport = c.options.transport
	}

	for _, wrapper := range c.options.transportWrappers {
		c.httpClient.Transport = wrapper(c.httpClient.Transport)
	}

	return c, nil
}
