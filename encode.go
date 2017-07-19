package ltsv

import (
	"bytes"
	"encoding"
	"fmt"
	"io"
	"reflect"
	"strconv"
	"strings"
	"sync"
)

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
	err := m[name]
	if e, ok := err.(*MarshalTypeError); ok {
		if e.err != nil {
			return e.err
		}
	}
	return m[name]
}

// An MarshalTypeError describes a LTSV value that was
// not appropriate for a value of a specific Go type.
type MarshalTypeError struct {
	Value string
	Type  reflect.Type
	key   string
	err   error
}

func (e *MarshalTypeError) Error() string {
	if e.err != nil {
		return e.err.Error()
	}
	return fmt.Sprintf("ltsv: failed to marshal type: %s, value: %s", e.Type.String(), e.Value)
}

var keyDelim = []byte{':'}
var valDelim = []byte{'\t'}

type fieldWriter func(w io.Writer, v reflect.Value) error

func makeStructWriter(v reflect.Value) fieldWriter {
	t := v.Type()
	n := t.NumField()

	writers := make([]fieldWriter, n)
	for i := 0; i < n; i++ {
		ft := t.Field(i)
		tag := ft.Tag.Get("ltsv")
		tags := strings.Split(tag, ",")
		key := tags[0]
		if key == "-" {
			continue
		}
		if key == "" {
			key = strings.ToLower(ft.Name)
		}
		kind := ft.Type.Kind()

		dereference := false
		if kind == reflect.Ptr {
			kind = ft.Type.Elem().Kind()
			dereference = true
		}

		var writer fieldWriter
		switch kind {
		case reflect.String:
			writer = makeStringWriter(key)
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			writer = makeIntWriter(key)
		case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
			writer = makeUintWriter(key)
		case reflect.Float32, reflect.Float64:
			writer = makeFloatWriter(key)
		default:
			dereference = false
			writer = makeInterfaceWriter(key)
		}
		if i > 0 {
			writer = withDelimWriter(writer)
		}
		if dereference {
			writer = elemWriter(writer)
		}
		writers[i] = writer
	}

	return fieldWriter(func(w io.Writer, v reflect.Value) error {
		errs := make(MarshalError)
		err := writers[0](w, v.Field(0))
		if err != nil {
			if e, ok := err.(*MarshalTypeError); ok {
				errs[e.key] = e
			}
		}

		for i, wr := range writers[1:] {
			if wr == nil {
				continue
			}
			err := wr(w, v.Field(i+1))
			if err != nil {
				if e, ok := err.(*MarshalTypeError); ok {
					errs[e.key] = e
				}
			}
		}
		if len(errs) > 0 {
			return errs
		}
		return nil
	})
}

func withDelimWriter(writer fieldWriter) fieldWriter {
	return fieldWriter(func(w io.Writer, v reflect.Value) error {
		w.Write(valDelim)
		return writer(w, v)
	})
}

func elemWriter(writer fieldWriter) fieldWriter {
	return fieldWriter(func(w io.Writer, v reflect.Value) error {
		if v.IsNil() {
			return nil
		}
		return writer(w, v.Elem())
	})
}

func writeField(w io.Writer, key, value string) {
	io.WriteString(w, key)
	w.Write(keyDelim)
	io.WriteString(w, value)
}

func makeStringWriter(key string) fieldWriter {
	return fieldWriter(func(w io.Writer, v reflect.Value) error {
		writeField(w, key, v.String())
		return nil
	})
}

func makeIntWriter(key string) fieldWriter {
	return fieldWriter(func(w io.Writer, v reflect.Value) error {
		writeField(w, key, strconv.FormatInt(v.Int(), 10))
		return nil
	})
}

func makeUintWriter(key string) fieldWriter {
	return fieldWriter(func(w io.Writer, v reflect.Value) error {
		writeField(w, key, strconv.FormatUint(v.Uint(), 10))
		return nil
	})
}

func makeFloatWriter(key string) fieldWriter {
	return fieldWriter(func(w io.Writer, v reflect.Value) error {
		writeField(w, key, strconv.FormatFloat(v.Float(), 'f', -1, v.Type().Bits()))
		return nil
	})
}

func makeInterfaceWriter(key string) fieldWriter {
	return fieldWriter(func(w io.Writer, v reflect.Value) error {
		if !v.CanInterface() {
			return &MarshalTypeError{key: key, Type: v.Type(), Value: v.String()}
		}

		switch u := v.Interface().(type) {
		case encoding.TextMarshaler:
			b, err := u.MarshalText()
			if err != nil {
				return &MarshalTypeError{key: key, Type: v.Type(), Value: v.String(), err: err}
			}
			io.WriteString(w, key)
			w.Write(keyDelim)
			w.Write(b)
			return nil
		default:
			return &MarshalTypeError{key: key, Type: v.Type(), Value: v.String()}
		}
	})
}

type writerCache struct {
	cache map[reflect.Type]fieldWriter
	sync.RWMutex
}

func (c *writerCache) Get(v reflect.Value) fieldWriter {
	c.RLock()
	t := v.Type()
	if v, ok := c.cache[t]; ok {
		c.RUnlock()
		return v
	}
	c.RUnlock()
	writer := makeStructWriter(v)

	c.Lock()
	c.cache[t] = writer
	c.Unlock()

	return writer
}

var cache = &writerCache{
	cache: make(map[reflect.Type]fieldWriter),
}

func marshalMapTo(w io.Writer, m map[string]string) error {
	first := true
	for k, v := range m {
		if !first {
			w.Write(valDelim)
		}
		first = false
		writeField(w, k, v)
	}
	return nil
}

func marshalStructTo(w io.Writer, rv reflect.Value) error {
	writer := cache.Get(rv)
	return writer(w, rv)
}

// MarshalTo writes the LTSV encoding of v into w.
// Be aware that the writing into w is not thread safe.
func MarshalTo(w io.Writer, v interface{}) error {
	rv := reflect.ValueOf(v)
	if rv.Kind() == reflect.Ptr {
		rv = rv.Elem()
	}

	var err error
	switch rv.Kind() {
	case reflect.Map:
		if m, ok := v.(map[string]string); ok {
			err = marshalMapTo(w, m)
			break
		}
		err = fmt.Errorf("not a map[string]string")
	case reflect.Struct:
		err = marshalStructTo(w, rv)
	default:
		err = fmt.Errorf("not a struct/map: %v", v)
	}
	return err
}

// Marshal returns the LTSV encoding of v
func Marshal(v interface{}) ([]byte, error) {
	w := bytes.NewBuffer(nil)
	err := MarshalTo(w, v)
	return w.Bytes(), err
}
