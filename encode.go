package ltsv

import (
	"fmt"
	"reflect"
	"strings"
)

type Marshaler interface {
	MarshalLTSV([]byte) error
}

// MarshalError is an error type for Marshal()
type MarshalError map[string]error

func (m MarshalError) Error() string {
	if len(m) == 0 {
		return "(no error)"
	}

	ee := make([]string, 0, len(m))
	for name, err := range m {
		ee = append(ee, fmt.Sprintf("field %q: %s", name, err))
	}

	return strings.Join(ee, "\n")
}

// OfField returns the error correspoinding to a given field
func (m MarshalError) OfField(name string) error {
	return m[name]
}

// An MarshalTypeError describes a LTSV value that was
// not appropriate for a value of a specific Go type.
type MarshalTypeError struct {
	Value string
	Type  reflect.Type
}

func (e *MarshalTypeError) Error() string {
	return "ltsv: cannot marshal " + e.Value + " into Go value of type " + e.Type.String()
}

func Marshal(v interface{}) ([]byte, error) {
	rv := reflect.ValueOf(v)
	if rv.Kind() == reflect.Ptr {
		rv = rv.Elem()
	}
	if rv.Kind() != reflect.Struct && rv.Kind() != reflect.Map {
		return nil, fmt.Errorf("not a struct/map: %v", v)
	}

	if rv.Kind() == reflect.Map {
		kt := rv.Type().Key()
		vt := rv.Type().Elem()
		if kt.Kind() != reflect.String || vt.Kind() != reflect.String {
			return nil, fmt.Errorf("not a map[string]string")
		}

		mKeys := rv.MapKeys()
		arr := make([]string, len(mKeys), len(mKeys))
		for i, k := range mKeys {
			arr[i] = k.String() + ":" + rv.MapIndex(k).String()
		}
		return []byte(strings.Join(arr, "\t")), nil
	}
	return nil, nil
}
