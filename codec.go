package apigen

import (
	"encoding/json"
	"fmt"
	"io"
	"strings"

	"github.com/morikuni/failure"
)

type Decoder interface {
	Decode(io.Reader) (*_struct, error)
}

type JSONDecoder struct{}

func (d *JSONDecoder) Decode(r io.Reader) (*_struct, error) {
	v := make(map[string]interface{})
	if err := json.NewDecoder(r).Decode(&v); err != nil {
		return nil, failure.Wrap(err)
	}

	return decodeJSONObject(v), nil
}

func decodeJSONObject(o map[string]interface{}) *_struct {
	var s _struct
	for k, v := range o {
		field := field{
			name:  k,
			_type: detectJSONType(v),
		}

		switch field._type {
		case typeStruct:
			field.value = &definedStruct{
				name:    strings.Title(k),
				_struct: decodeJSONObject(v.(map[string]interface{})),
			}
		case typeSlice:
			field.value = decodeJSONArray(v.([]interface{}))
		case typeBool:
			field.value = v
		case typeString:
			field.value = v
		case typeFloat64:
			field.value = v
		}

		s.fields = append(s.fields, field)
	}

	return &s
}

func decodeJSONArray(array []interface{}) *slice {
	if len(array) == 0 {
		return nil
	}

	slice := &slice{
		_type: detectJSONType(array[0]),
		elems: make([]interface{}, 0, len(array)),
	}

	for _, v := range array {
		switch slice._type {
		case typeStruct:
			slice.elems = append(slice.elems, decodeJSONObject(v.(map[string]interface{})))
		case typeSlice:
			slice.elems = append(slice.elems, decodeJSONArray(v.([]interface{})))
		case typeBool:
			slice.elems = append(slice.elems, v)
		case typeString:
			slice.elems = append(slice.elems, v)
		case typeFloat64:
			slice.elems = append(slice.elems, v)
		}
	}

	return slice
}

func detectJSONType(v interface{}) _type {
	switch cv := v.(type) {
	case map[string]interface{}:
		return typeStruct
	case []interface{}:
		return typeSlice
	case bool:
		return typeBool
	case string:
		return typeString
	case float64:
		return typeFloat64
	default:
		panic(fmt.Sprintf("unknown type: %T", cv))
	}
}
