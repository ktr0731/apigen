package apigen

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/url"
	"sort"
	"strings"

	"github.com/morikuni/failure"
)

type jsonType int

const (
	jsonTypeNull jsonType = iota
	jsonTypeBool
	jsonTypeNumber
	jsonTypeString
	jsonTypeArray
	jsonTypeObject
)

type decoder interface {
	Decode(io.Reader) (_type, error)
}

type jsonDecoder struct{}

func (d *jsonDecoder) Decode(r io.Reader) (_type, error) {
	var v interface{}
	if err := json.NewDecoder(r).Decode(&v); err != nil {
		return nil, failure.Wrap(err)
	}

	switch v := v.(type) {
	case map[string]interface{}:
		return decodeJSONObject(v), nil
	case []interface{}:
		return decodeJSONArray(v), nil
	default:
		return nil, errors.New("unsupported top-level JSON type")
	}
}

func decodeJSONType(v interface{}) _type {
	switch detectJSONType(v) {
	case jsonTypeNull:
		return emptyIfaceType
	case jsonTypeObject:
		return decodeJSONObject(v.(map[string]interface{}))
	case jsonTypeArray:
		return decodeJSONArray(v.([]interface{}))
	case jsonTypeBool:
		return typeBool
	case jsonTypeString:
		return typeString
	case jsonTypeNumber:
		return typeFloat64
	default:
		panic("unreachable")
	}
}

func decodeJSONObject(o map[string]interface{}) *structType {
	var s structType
	for k, v := range o {
		key := public(k)
		field := &structField{
			name: key,
			tags: map[string][]string{"json": {k, "omitempty"}},
		}

		switch detectJSONType(v) {
		case jsonTypeNull:
			field._type = emptyIfaceType
		case jsonTypeObject:
			field._type = decodeJSONObject(v.(map[string]interface{}))
		case jsonTypeArray:
			field._type = decodeJSONArray(v.([]interface{}))
		case jsonTypeBool:
			field._type = typeBool
		case jsonTypeString:
			field._type = typeString
		case jsonTypeNumber:
			field._type = typeFloat64
		}

		s.fields = append(s.fields, field)
	}

	sort.Slice(s.fields, func(i, j int) bool {
		return s.fields[i].name < s.fields[j].name
	})

	return &s
}

func decodeJSONArray(arr []interface{}) *sliceType {
	if len(arr) == 0 {
		return &sliceType{elemType: emptyIfaceType}
	}

	return &sliceType{elemType: decodeJSONType(arr[0])}
}

func detectJSONType(v interface{}) jsonType {
	switch cv := v.(type) {
	case nil:
		return jsonTypeNull
	case map[string]interface{}:
		return jsonTypeObject
	case []interface{}:
		return jsonTypeArray
	case bool:
		return jsonTypeBool
	case string:
		return jsonTypeString
	case float64:
		return jsonTypeNumber
	default:
		panic(fmt.Sprintf("unknown type: %T", cv))
	}
}

func structFromPathParams(h string, u *url.URL) *structType {
	if h == "" {
		return &structType{}
	}

	var (
		t           structType
		left        int
		replaceArgs []string
	)

	runes := []rune(h)
	for i, r := range runes {
		switch r {
		case '{':
			left = i
		case '}':
			k := string(runes[left+1 : i]) // Param name without '{' and '}'.
			t.fields = append(t.fields, &structField{
				name:  public(k), // TODO: Support snake case.
				_type: typeString,
			})
			replaceArgs = append(replaceArgs, string(runes[left:i+1]), "%s")
		}
	}

	u.Path = strings.NewReplacer(replaceArgs...).Replace(h)

	return &t
}

func structFromQuery(q url.Values) *structType {
	var s structType
	for k, v := range q {
		field := &structField{
			name: public(k),
			meta: map[string]string{"key": k},
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
