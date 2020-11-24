package apigen

import (
	"fmt"
	"go/format"
	"io"
	"strconv"
	"strings"
	"sync"

	"github.com/morikuni/failure"
)

type method struct {
	name string
	req  *Request
	res  *Response
}

type Generator struct {
	writer  io.Writer
	b       strings.Builder
	err     error
	errOnce sync.Once
	methods []*method
}

func NewGenerator(w io.Writer) *Generator {
	return &Generator{writer: w}
}

func (g *Generator) Add(name string, req *Request, res *Response) {
	g.methods = append(g.methods, &method{name: name, req: req, res: res})
}

func (g *Generator) Generate(s *_struct) error {
	g._package("main")
	g._import("fmt")
	// for _, m := range g.methods {
	g.typeStruct("Method", s)
	// }

	out, err := g.gen()
	if err != nil {
		return failure.Wrap(err)
	}

	b, err := format.Source([]byte(out))
	if err != nil {
		return failure.Wrap(err)
	}

	if _, err := g.writer.Write(b); err != nil {
		return failure.Wrap(err)
	}

	return nil
}

func (g *Generator) _package(name string) {
	g.wf("package %s", name)
}

func (g *Generator) _import(paths ...string) {
	g.w("import (")
	for _, path := range paths {
		g.w(strconv.Quote(path))
	}
	g.w(")")
}

func (g *Generator) typeStruct(name string, s *_struct) {
	g.wf("type %s struct {", name)
	for _, f := range s.fields {
		g.wf("%s %s", f.name, f._type)
	}
	g.w("}")
}

func (g *Generator) w(s string) {
	if g.err != nil {
		return
	}

	_, err := io.WriteString(&g.b, s+"\n")
	g.error(err)
}

func (g *Generator) wf(f string, a ...interface{}) {
	if g.err != nil {
		return
	}

	_, err := fmt.Fprintf(&g.b, f+"\n", a...)
	g.error(err)
}

func (g *Generator) error(err error) {
	if err == nil {
		return
	}
	g.errOnce.Do(func() {
		g.err = err
	})
}

func (g *Generator) gen() (string, error) {
	if g.err != nil {
		return "", g.err
	}
	return g.b.String(), nil
}
