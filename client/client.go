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
	headers    http.Header
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

	for k, v := range c.headers {
		for _, vv := range v {
			hreq.Header.Add(k, vv)
		}
	}

	hres, err := c.httpClient.Do(hreq)
	if err != nil {
		return err
	}
	defer hres.Body.Close()

	// TODO: Support another codecs.
	return json.NewDecoder(hres.Body).Decode(&res)
}

func reqToRawQuery(v interface{}) string {
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
