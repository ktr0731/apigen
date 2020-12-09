package client

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/url"
)

var defaultClient = Client{
	httpClient: http.DefaultClient,
}

type Client struct {
	httpClient *http.Client
	ints       []Interceptor
}

func New(opts ...Option) *Client {
	c := defaultClient

	for _, o := range opts {
		o(&c)
	}

	return &c
}

func (c *Client) Do(
	ctx context.Context,
	method string,
	url *url.URL,
	req, res interface{},
) error {
	var body bytes.Buffer
	if req != nil {
		if err := json.NewEncoder(&body).Encode(&req); err != nil {
			return err
		}
	}

	hreq, err := http.NewRequestWithContext(ctx, method, url.String(), &body)
	if err != nil {
		return err
	}

	chainedHandler := func(ctx context.Context, req *http.Request) (*http.Response, error) {
		return c.httpClient.Do(req.WithContext(ctx))
	}
	chainer := func(i Interceptor, h Handler) Handler {
		return func(ctx context.Context, req *http.Request) (*http.Response, error) {
			return i(ctx, req.WithContext(ctx), h)
		}
	}

	for i := len(c.ints) - 1; i >= 0; i-- {
		chainedHandler = chainer(c.ints[i], chainedHandler)
	}

	hres, err := chainedHandler(ctx, hreq)
	if err != nil {
		return err
	}
	defer hres.Body.Close()

	// TODO: Support another codecs.
	return json.NewDecoder(hres.Body).Decode(&res)
}
