package client_test

import (
	"context"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
	"testing"

	"github.com/ktr0731/apigen/client"
)

type req struct {
	Foo string `json:"foo"`
}

type res struct {
	Bar int `json:"bar"`
}

type transport struct {
	roundTripFn func(*http.Request) (*http.Response, error)
}

func (t *transport) RoundTrip(r *http.Request) (*http.Response, error) {
	return t.roundTripFn(r)
}

func TestClient(t *testing.T) {
	t.Parallel()

	cases := map[string]struct {
		interceptors []client.Interceptor
		roundTripFn  func(*http.Request) (*http.Response, error)
		wantErr      error
	}{
		"ok": {
			roundTripFn: func(*http.Request) (*http.Response, error) {
				return &http.Response{
					Body: ioutil.NopCloser(strings.NewReader(`{ "bar": 100 }`)),
				}, nil
			},
		},
		"with HeaderInterceptor": {
			interceptors: []client.Interceptor{
				client.HeaderInterceptor(http.Header{"foo": []string{"bar"}}),
				func(_ context.Context, req *http.Request, _ client.Handler) (*http.Response, error) {
					if want, got := "bar", req.Header.Get("foo"); want != got {
						return nil, fmt.Errorf("want header value %s, but got %s", want, got)
					}

					return &http.Response{
						Body: ioutil.NopCloser(strings.NewReader(`{ "bar": 100 }`)),
					}, nil
				},
			},
		},
		"invalid response body": {
			roundTripFn: func(*http.Request) (*http.Response, error) {
				return &http.Response{
					Body: ioutil.NopCloser(strings.NewReader("")),
				}, nil
			},
			wantErr: io.EOF,
		},
		"second interceptor returns io.ErrUnexpectedEOF": {
			interceptors: []client.Interceptor{
				func(ctx context.Context, req *http.Request, h client.Handler) (*http.Response, error) {
					return h(ctx, req)
				},
				func(context.Context, *http.Request, client.Handler) (*http.Response, error) {
					return nil, io.ErrUnexpectedEOF
				},
			},
			wantErr: io.ErrUnexpectedEOF,
		},
	}

	for name, c := range cases {
		c := c

		t.Run(name, func(t *testing.T) {
			t.Parallel()

			cli := client.New(
				client.WithHTTPClient(&http.Client{Transport: &transport{roundTripFn: c.roundTripFn}}),
				client.WithInterceptors(c.interceptors...),
			)

			u, err := url.Parse("http://example.com")
			if err != nil {
				t.Fatalf("url.Parse should not return an error, but got '%s'", err)
			}

			err = cli.Do(context.Background(), "GET", u, &req{}, &res{})
			if c.wantErr != nil {
				if !errors.Is(err, c.wantErr) {
					t.Errorf("want error '%s', but got '%s'", c.wantErr, err)
				}

				return
			}
			if err != nil {
				t.Fatalf("should not return an error, but got '%s'", err)
			}
		})
	}
}
