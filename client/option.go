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

func WithHeaders(h http.Header) Option {
	return func(client *Client) {
		client.headers = h
	}
}
