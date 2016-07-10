package ltsv

import (
	"encoding"
	"fmt"
	"reflect"
	"strconv"
	"strings"
)

// Marshaler is the interface inmpemented by types that can marshal themselves
type Marshaler interface {
	MarshalLTSV() ([]byte, error)
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
	return "ltsv: cannot marshal Go value " + e.Value + " of type " + e.Type.String() + " into ltsv"
}

// Marshal returns the LTSV encoding of v
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

	t := rv.Type()
	numField := t.NumField()
	arr := make([]string, 0, numField)
	errs := MarshalError{}
	for i := 0; i < numField; i++ {
		ft := t.Field(i)
		fv := rv.Field(i)
		tag := ft.Tag.Get("ltsv")
		tags := strings.Split(tag, ",")
		key := tags[0]
		if key == "-" {
			continue
		}
		if key == "" {
			key = strings.ToLower(ft.Name)
		}

		switch fv.Kind() {
		case reflect.String:
			arr = append(arr, key+":"+fv.String())
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			arr = append(arr, key+":"+strconv.FormatInt(fv.Int(), 10))
		case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
			arr = append(arr, key+":"+strconv.FormatUint(fv.Uint(), 10))
		case reflect.Float32, reflect.Float64:
			arr = append(arr, key+":"+strconv.FormatFloat(fv.Float(), 'f', -1, fv.Type().Bits()))
		case reflect.Interface:
			if u, ok := fv.Interface().(Marshaler); ok {
				buf, err := u.MarshalLTSV()
				if err != nil {
					errs[ft.Name] = err
				} else {
					arr = append(arr, key+":"+string(buf))
				}
				continue
			}
			if u, ok := fv.Interface().(encoding.TextMarshaler); ok {
				buf, err := u.MarshalText()
				if err != nil {
					errs[ft.Name] = err
				} else {
					arr = append(arr, key+":"+string(buf))
				}
				continue
			}
			fallthrough
		default:
			errs[ft.Name] = &MarshalTypeError{fv.String(), fv.Type()}
		}
	}
	if len(errs) < 1 {
		return []byte(strings.Join(arr, "\t")), nil
	}
	return nil, errs
}
