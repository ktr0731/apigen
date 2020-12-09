package apigen

import (
	"context"
	"io"
	"net/http"
	"os"
	"strings"

	"github.com/morikuni/failure"
	"golang.org/x/sync/errgroup"
)

// Generate generates client interfaces, types for requests and responses based on the passed API definitions.
// If the API definition is invalid, Generate returns an error wrapping ErrInvalidDefinition.
func Generate(ctx context.Context, def *Definition, opts ...Option) error {
	return newRunner(opts...).run(ctx, def)
}

type runner struct {
	client  *http.Client
	decoder decoder
	writer  io.Writer
	pkg     string
}

func newRunner(opts ...Option) *runner {
	r := &runner{
		client:  http.DefaultClient,
		decoder: &jsonDecoder{},
		writer:  os.Stdout,
		pkg:     "main",
	}
	for _, o := range opts {
		o(r)
	}

	return r
}

func (r *runner) run(ctx context.Context, def *Definition) error {
	if err := def.validate(); err != nil {
		return failure.Wrap(err)
	}

	eg, cctx := errgroup.WithContext(ctx)
	for service, methods := range def.Services {
		eg.Go(func() error {
			return r.processService(cctx, service, methods)
		})
	}

	return eg.Wait()
}

func (r *runner) processService(ctx context.Context, service string, methods []*Method) error {
	eg, cctx := errgroup.WithContext(ctx)
	gen := newGenerator(r.writer)

	for _, m := range methods {
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
				// TODO: Call before calling API.
				if m.ParamHint != "" {
					methReq = structFromPathParams(m.ParamHint, req.URL)
				} else {
					methReq = structFromQuery(req.URL.Query())
				}
			default:
				panic("not implemnted yet")
			}

			u := req.URL
			u.RawQuery = ""
			gen.addMethod(service+"Client", &method{
				name:   m.Name,
				method: req.Method,
				url:    strings.ReplaceAll(u.String(), "%25s", "%s"), // Replace URL-encoded '%s'.
				req:    methReq,
				res:    methRes,
			})

			return nil
		})
	}

	if err := eg.Wait(); err != nil {
		return err
	}

	return gen.generate(r.pkg)
}
