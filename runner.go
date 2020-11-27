package apigen

import (
	"context"
	"net/http"

	"github.com/k0kubun/pp"
	"github.com/morikuni/failure"
)

type Option func(*Runner)

type Runner struct {
	client  *http.Client
	decoder Decoder
}

var defaultRunner = Runner{
	client:  http.DefaultClient,
	decoder: &JSONDecoder{},
}

func NewRunner(opts ...Option) *Runner {
	r := defaultRunner
	for _, o := range opts {
		o(&r)
	}

	return &r
}

type Request struct {
	name    string
	_struct *_struct
}

type Response struct {
	name    string
	_struct *_struct
}

func (r *Runner) Run(ctx context.Context, req *http.Request) (*Request, *Response, error) {
	res, err := r.client.Do(req)
	if err != nil {
		return nil, nil, failure.Wrap(err)
	}
	defer res.Body.Close()

	out, err := r.decoder.Decode(res.Body)
	if err != nil {
		return nil, nil, failure.Wrap(err)
	}

	pp.Println(out)

	return nil, &Response{name: "Response", _struct: out}, nil
}
