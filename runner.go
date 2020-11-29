package apigen

import (
	"context"
	"net/http"
	"os"

	"github.com/morikuni/failure"
	"golang.org/x/sync/errgroup"
)

func Generate(ctx context.Context, def *Definition, opts ...Option) error {
	return newRunner(opts...).Run(ctx, def)
}

type runner struct {
	client  *http.Client
	decoder decoder
}

var defaultRunner = runner{
	client:  http.DefaultClient,
	decoder: &jsonDecoder{},
}

func newRunner(opts ...Option) *runner {
	r := defaultRunner
	for _, o := range opts {
		o(&r)
	}

	return &r
}

func (r *runner) Run(ctx context.Context, def *Definition) error {
	if err := def.validate(); err != nil {
		return failure.Wrap(err)
	}

	gen := newGenerator(os.Stdout)

	eg, cctx := errgroup.WithContext(ctx)

	for _, m := range def.Methods {
		req, err := m.Request(cctx)
		if err != nil {
			return failure.Wrap(err)
		}

		m := m

		eg.Go(func() error {
			res, err := r.client.Do(req)
			if err != nil {
				return failure.Wrap(err)
			}
			defer res.Body.Close()

			methRes, err := r.decoder.Decode(res.Body)
			if err != nil {
				return failure.Wrap(err)
			}

			var methReq *structType
			switch req.Method {
			case http.MethodGet:
				methReq = structFromQuery(req.URL.Query())
			}

			for k, v := range req.Header {
				for _, vv := range v {
					gen.addHeader(k, vv)
				}
			}

			u := req.URL
			u.RawQuery = ""
			gen.addMethod(m.Service+"Client", &method{
				Name:   m.Method,
				method: req.Method,
				url:    u.String(),
				req:    methReq,
				res:    methRes,
			})

			return nil
		})
	}

	if err := eg.Wait(); err != nil {
		return err
	}

	return gen.generate()
}
