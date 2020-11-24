package apigen

type _type string

const (
	typeBool   _type = "bool"
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

	typeString _type = "string"
	typeFloat  _type = "float"

	typeStruct _type = "struct"
	typeArray  _type = "array"
	typeSlice  _type = "slice"
)

type field struct {
	name  string
	_type _type
}

type _struct struct {
	fields []field
}
