package apigen

import (
	"fmt"
	"path"
	"strings"

	"github.com/achiku/varfmt"
)

var (
	typeBool    _type = basicType("bool")
	typeString  _type = basicType("string")
	typeFloat64 _type = basicType("float64")
)

type _type interface {
	String() string
	isBasic() bool
}

type basicType string

func (t basicType) String() string { return string(t) }
func (t basicType) isBasic() bool  { return true }

type sliceType struct {
	elemType _type
}

func (t *sliceType) isBasic() bool  { return false }
func (t *sliceType) String() string { return fmt.Sprintf("[]%s", t.elemType.String()) }

type structType struct {
	fields []*structField
}

func (t *structType) isBasic() bool { return false }
func (t *structType) String() string {
	if len(t.fields) == 0 {
		return "struct {}"
	}

	s := "struct {\n"
	for i := range t.fields {
		s += fmt.Sprintf("%s\n", t.fields[i].String())
	}
	s += "}"

	return s
}

type structField struct {
	name  string // Means embedded field if name is empty.
	_type _type
	tags  map[string][]string
}

func (f *structField) String() string {
	s := f._type.String()
	if f.name != "" {
		s = fmt.Sprintf("%s %s", f.name, s)
	}

	if len(f.tags) != 0 {
		tags := make([]string, 0, len(f.tags))
		for k, v := range f.tags {
			tags = append(tags, fmt.Sprintf(`%s:"%s"`, k, strings.Join(v, ",")))
		}
		s += fmt.Sprintf(" `%s`", strings.Join(tags, " "))
	}

	return s
}

var emptyIfaceType = &_emptyIfaceType{}

type _emptyIfaceType struct{}

func (t *_emptyIfaceType) isBasic() bool  { return false }
func (t *_emptyIfaceType) String() string { return "interface{}" }

type definedType struct {
	pkg     string
	name    string
	pointer bool
	_type   _type // Nil if pkg is empty (declared by another package).
}

func (t *definedType) isBasic() bool { return false }
func (t *definedType) String() string {
	var s string
	if t.pkg != "" {
		s = fmt.Sprintf("%s.%s", path.Base(t.pkg), t.name)
	} else {
		s = t.name
	}
	if t.pointer {
		s = "*" + s
	}

	return s
}

func public(s string) string { return varfmt.PublicVarName(s) }
