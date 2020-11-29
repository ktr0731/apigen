package client

import (
	"net/http"
)

type Option func(*Client)

func WithHTTPClient(c *http.Client) Option {
	return func(client *Client) {
		client.httpClient = c
	}
}

func WithInterceptors(ints ...Interceptor) Option {
	return func(client *Client) {
		client.ints = append(client.ints, ints...)
	}
}
