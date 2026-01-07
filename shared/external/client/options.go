package client

import "net/http"

type clientOptions struct {
	httpClient        *http.Client
	transport         http.RoundTripper
	transportWrappers []TransportWrapper
}

// ClientOption defines the properties of a client option.
type ClientOption func(*clientOptions) error

// WithHTTPClient sets the http client of the base client.
func WithHTTPClient(httpClient *http.Client) ClientOption {
	return func(c *clientOptions) error {
		c.httpClient = httpClient

		return nil
	}
}

// WithTransport sets the transport of the http client of the base client.
func WithTransport(transport http.RoundTripper) ClientOption {
	return func(c *clientOptions) error {
		c.transport = transport

		return nil
	}
}

// TransportWrapper defines the properties of a transport wrapper.
type TransportWrapper func(http.RoundTripper) http.RoundTripper

// WithTransportWrappers wraps the transport of the http client of the base client.
func WithTransportWrappers(wrappers ...TransportWrapper) ClientOption {
	return func(c *clientOptions) error {
		c.transportWrappers = append(c.transportWrappers, wrappers...)

		return nil
	}
}
