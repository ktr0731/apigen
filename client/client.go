package client

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/url"
	"reflect"
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
	var body io.Reader
	switch method {
	case http.MethodGet:
		url.RawQuery = reqToRawQuery(req)
	case http.MethodPost:
		panic("not implemented yet")
	}

	hreq, err := http.NewRequestWithContext(ctx, method, url.String(), body)
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

func reqToRawQuery(v interface{}) string {
	if v == nil {
		return ""
	}

	rv := indirect(reflect.ValueOf(v))
	rt := rv.Type()

	vals := make(url.Values)
	for i := 0; i < rt.NumField(); i++ {
		k := rt.Field(i).Tag.Get("name")
		v := rv.Field(i).Interface().(string)
		vals.Add(k, v)
	}

	return vals.Encode()
}

func indirect(rv reflect.Value) reflect.Value {
	if rv.Type().Kind() != reflect.Ptr {
		return rv
	}

	return indirect(reflect.Indirect(rv))
}
