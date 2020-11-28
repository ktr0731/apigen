package apigen

import (
	"context"
	"net/http"
	"path"

	"github.com/iancoleman/strcase"
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

type Method struct {
	Name string
	req  *structType
	res  *structType
}

func (r *Runner) Run(ctx context.Context, req *http.Request) (*Method, error) {
	res, err := r.client.Do(req)
	if err != nil {
		return nil, failure.Wrap(err)
	}
	defer res.Body.Close()

	methRes, err := r.decoder.Decode(res.Body)
	if err != nil {
		return nil, failure.Wrap(err)
	}

	var methReq *structType
	switch req.Method {
	case http.MethodGet:
		methReq = structFromQuery(req.URL.Query())
	}

	pp.Println(methRes)

	return &Method{
		Name: strcase.ToCamel(public(path.Base(req.URL.Path))),
		req:  methReq,
		res:  methRes,
	}, nil
}
