package apigen

import (
	"context"
	"net/http"
	"net/url"
	"path"

	"github.com/iancoleman/strcase"
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
		methReq = reqStructFromQuery(req.URL.Query())
	}

	return &Method{
		Name: strcase.ToCamel(public(path.Base(req.URL.Path))),
		req:  methReq,
		res:  methRes,
	}, nil
}

func reqStructFromQuery(q url.Values) *structType {
	var s structType
	for k, v := range q {
		field := &structField{
			name: public(k),
		}
		if len(v) == 1 {
			field._type = typeString
		} else {
			field._type = &sliceType{elemType: typeString}
		}
		s.fields = append(s.fields, field)
	}
	return &s
}
