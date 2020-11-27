package apigen

type _type string

func (t _type) isBasic() bool {
	switch t {
	case typeBool, typeString, typeFloat64:
		return true
	default:
		return false
	}
}

const (
	typeBool   _type = "bool"
	typeString _type = "string"
	typeInt    _type = "int"
	typeInt8   _type = "int8"
	typeInt16  _type = "int16"
	typeInt32  _type = "int32"
	typeInt64  _type = "int64"
	typeUint   _type = "uint"
	typeUint8  _type = "uint8"
	typeUint16 _type = "uint16"
	typeUint32 _type = "uint32"
	typeUint64 _type = "uint64"

	typeFloat   _type = "float"
	typeFloat64 _type = "float64"

	typeStruct _type = "struct"
	typeArray  _type = "array"
	typeSlice  _type = "slice"
)

type field struct {
	name  string
	_type _type
	value interface{}
}

type definedStruct struct {
	name    string
	_struct *_struct
}

type _struct struct {
	fields []*field
}

type slice struct {
	_type _type
	elems []interface{}
}
