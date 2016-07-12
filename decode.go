package ltsv

import (
	"encoding"
	"fmt"
	"reflect"
	"strconv"
	"strings"
)

// UnmarshalError is an error type for Unmarshal()
type UnmarshalError map[string]error

func (m UnmarshalError) Error() string {
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
func (m UnmarshalError) OfField(name string) error {
	return m[name]
}

// An UnmarshalTypeError describes a LTSV value that was
// not appropriate for a value of a specific Go type.
type UnmarshalTypeError struct {
	Value string
	Type  reflect.Type
}

func (e *UnmarshalTypeError) Error() string {
	return "ltsv: cannot unmarshal " + e.Value + " into Go value of type " + e.Type.String()
}

type ltsvMap map[string]string

func data2map(data []byte) (ltsvMap, error) {
	d := string(data)
	fields := strings.Split(d, "\t")
	l := ltsvMap{}
	for _, v := range fields {
		kv := strings.SplitN(strings.TrimSpace(v), ":", 2)
		if len(kv) != 2 {
			return nil, fmt.Errorf("not a ltsv: %s", d)
		}
		l[strings.TrimSpace(kv[0])] = strings.TrimSpace(kv[1])
	}
	return l, nil
}

// Unmarshal parses the LTSV-encoded data and stores the result
// in the value pointed to by v.
func Unmarshal(data []byte, v interface{}) error {
	rv := reflect.ValueOf(v)
	if rv.Kind() != reflect.Ptr {
		return fmt.Errorf("not a pointer: %v", v)
	}

	rv = rv.Elem()
	if rv.Kind() != reflect.Struct && rv.Kind() != reflect.Map {
		return fmt.Errorf("not a pointer to a struct/map: %v", v)
	}

	l, err := data2map(data)
	if err != nil {
		return err
	}
	if rv.Kind() == reflect.Map {
		kt := rv.Type().Key()
		vt := rv.Type().Elem()
		if kt.Kind() != reflect.String || vt.Kind() != reflect.String {
			return fmt.Errorf("not a map[string]string")
		}
		for k, v := range l {
			kv := reflect.ValueOf(k).Convert(kt)
			vv := reflect.ValueOf(v).Convert(vt)
			rv.SetMapIndex(kv, vv)
		}
		return nil
	}

	t := rv.Type()
	errs := UnmarshalError{}
	for i := 0; i < t.NumField(); i++ {
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
		s, ok := l[key]
		if !ok {
			continue
		}

		if fv.Kind() == reflect.Ptr {
			if fv.IsNil() {
				fv.Set(reflect.New(fv.Type().Elem()))
			}
			fv = fv.Elem()
		}
		if !fv.CanSet() {
			continue
		}

		switch fv.Kind() {
		case reflect.String:
			fv.SetString(s)
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			i, err := strconv.ParseInt(s, 10, 64)
			if err != nil || fv.OverflowInt(i) {
				errs[ft.Name] = &UnmarshalTypeError{"number " + s, fv.Type()}
				continue
			}
			fv.SetInt(i)
		case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
			i, err := strconv.ParseUint(s, 10, 64)
			if err != nil || fv.OverflowUint(i) {
				errs[ft.Name] = &UnmarshalTypeError{"number " + s, fv.Type()}
				continue
			}
			fv.SetUint(i)
		case reflect.Float32, reflect.Float64:
			n, err := strconv.ParseFloat(s, fv.Type().Bits())
			if err != nil || fv.OverflowFloat(n) {
				errs[ft.Name] = &UnmarshalTypeError{"number " + s, fv.Type()}
				continue
			}
			fv.SetFloat(n)
		case reflect.Interface:
			if tu, ok := fv.Interface().(encoding.TextUnmarshaler); ok {
				err := tu.UnmarshalText([]byte(s))
				if err != nil {
					errs[ft.Name] = err
				}
				continue
			}
			fallthrough
		default:
			errs[ft.Name] = &UnmarshalTypeError{s, fv.Type()}
		}
	}

	if len(errs) < 1 {
		return nil
	}
	return errs
}
