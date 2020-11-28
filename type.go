package apigen

import (
	"fmt"
	"path"
)

var (
	typeBool    _type = basicType("bool")
	typeString  _type = basicType("string")
	typeFloat64 _type = basicType("float64")
)

type _type interface {
	name() string
	isBasic() bool
}

type basicType string

func (t basicType) name() string  { return string(t) }
func (t basicType) isBasic() bool { return true }

type sliceType struct {
	elemType _type
}

func (t *sliceType) isBasic() bool { return false }
func (t *sliceType) name() string  { return fmt.Sprintf("[]%s", t.elemType.name()) }

type structType struct {
	fields []*structField
}

func (t *structType) isBasic() bool { return false }
func (t *structType) name() string {
	s := "struct {\n"
	for i := range t.fields {
		s += fmt.Sprintf("%s %s\n", t.fields[i].name, t.fields[i]._type.name())
	}
	s += "}\n"
	return s
}

type structField struct {
	name  string
	_type _type
}

var emptyIfaceType = &_emptyIfaceType{}

type _emptyIfaceType struct{}

func (t *_emptyIfaceType) isBasic() bool { return false }
func (t *_emptyIfaceType) name() string  { return "interface{}" }

type externalType struct{}

func (t *externalType) isBasic() bool { return false }
func (t *externalType) name() string  { return "interface{}" }

type definedType struct {
	pkg   string
	tName string
	_type _type // Nil if pkg is not empty (declared by another package).
}

func (t *definedType) isBasic() bool { return false }
func (t *definedType) name() string {
	if t.pkg != "" {
		return fmt.Sprintf("*%s.%s", path.Base(t.pkg), t.tName)
	}
	return fmt.Sprintf("*%s", t.tName)
}
