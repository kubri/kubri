package desc

import (
	"bytes"
	"encoding"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"reflect"
	"strconv"
	"sync"
	"time"
)

var encoders sync.Map //nolint:gochecknoglobals

type encoder func(io.Writer, reflect.Value) error

func Marshal(v any) ([]byte, error) {
	buf := &bytes.Buffer{}
	if err := NewEncoder(buf).Encode(v); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

type Encoder struct {
	w io.Writer
}

func NewEncoder(w io.Writer) *Encoder {
	return &Encoder{w}
}

func (e *Encoder) Encode(v any) error {
	if v == nil {
		return errors.New("unsupported type: nil")
	}

	val := reflect.ValueOf(v)
	typ := val.Type()

	if enc, ok := encoders.Load(typ); ok {
		return enc.(encoder)(e.w, val) //nolint:forcetypeassert
	}

	enc, err := newEncoder(typ)
	if err != nil {
		return err
	}
	encoders.Store(typ, enc)
	return enc(e.w, val)
}

//nolint:gocognit,funlen
func newEncoder(typ reflect.Type) (encoder, error) {
	encoders := []func(buf *bytes.Buffer, w io.Writer, v reflect.Value) error{}

	var isPointer int8
	for typ.Kind() == reflect.Pointer {
		typ = typ.Elem()
		isPointer++
	}
	if typ.Kind() != reflect.Struct {
		return nil, fmt.Errorf("unsupported type: %s", typ)
	}

	for i := range typ.NumField() {
		field := typ.Field(i)

		n := getFieldName(field)
		if n == "" {
			continue
		}
		name := atob(n)

		enc, err := newFieldEncoder(field.Type)
		if err != nil {
			return nil, err
		}
		encoders = append(encoders, func(buf *bytes.Buffer, w io.Writer, v reflect.Value) error {
			if err := enc(buf, v.Field(i)); err != nil {
				return err
			}
			if buf.Len() == 0 {
				return nil
			}
			if _, err = w.Write(percent); err != nil {
				return err
			}
			if _, err = w.Write(name); err != nil {
				return err
			}
			if _, err = w.Write(percent); err != nil {
				return err
			}
			if _, err = w.Write(nl); err != nil {
				return err
			}
			if _, err = io.Copy(w, buf); err != nil {
				return err
			}
			if _, err = w.Write(nl); err != nil {
				return err
			}
			_, err = w.Write(nl)
			return err
		})
	}

	return func(w io.Writer, v reflect.Value) error {
		for range isPointer {
			if v.IsNil() {
				return nil
			}
			v = v.Elem()
		}
		buf := bufPool.Get().(*bytes.Buffer) //nolint:forcetypeassert
		defer bufPool.Put(buf)
		for _, enc := range encoders {
			buf.Reset()
			if err := enc(buf, w, v); err != nil {
				return err
			}
		}
		return nil
	}, nil
}

//nolint:gochecknoglobals
var (
	stringerType  = reflect.TypeOf((*fmt.Stringer)(nil)).Elem()
	marshalerType = reflect.TypeOf((*encoding.TextMarshaler)(nil)).Elem()
)

func newFieldEncoder(typ reflect.Type) (encoder, error) {
	switch {
	case typ == dateType:
		return newDateEncoder(typ)
	case typ.Implements(marshalerType):
		return newMarshalerEncoder(typ)
	case typ.Implements(stringerType):
		return newStringerEncoder(typ)
	}

	switch typ.Kind() {
	case reflect.Ptr:
		return newPtrEncoder(typ)
	case reflect.Slice:
		return newSliceEncoder(typ)
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return newIntEncoder(typ)
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return newUintEncoder(typ)
	case reflect.Float32, reflect.Float64:
		return newFloatEncoder(typ)
	case reflect.String:
		return newStringEncoder(typ)
	case reflect.Array:
		if typ.Elem().Kind() == reflect.Uint8 {
			return newByteArrayEncoder(typ)
		}
	}
	return nil, fmt.Errorf("unsupported type: %s", typ)
}

func newStringerEncoder(typ reflect.Type) (encoder, error) {
	ptr := typ.Kind() == reflect.Pointer
	return func(w io.Writer, v reflect.Value) error {
		if ptr && v.IsNil() {
			return nil
		}
		s := v.Interface().(fmt.Stringer).String() //nolint:forcetypeassert
		if len(s) == 0 {
			return nil
		}
		_, err := w.Write(atob(s))
		return err
	}, nil
}

func newMarshalerEncoder(typ reflect.Type) (encoder, error) {
	ptr := typ.Kind() == reflect.Pointer
	return func(w io.Writer, v reflect.Value) error {
		if ptr && v.IsNil() {
			return nil
		}
		b, err := v.Interface().(encoding.TextMarshaler).MarshalText() //nolint:forcetypeassert
		if err != nil {
			return err
		}
		if len(b) == 0 {
			return nil
		}
		_, err = w.Write(b)
		return err
	}, nil
}

func newPtrEncoder(typ reflect.Type) (encoder, error) {
	typ = typ.Elem()
	enc, err := newFieldEncoder(typ)
	if err != nil {
		return nil, err
	}

	return func(w io.Writer, v reflect.Value) error {
		if v.IsNil() {
			return nil
		}
		return enc(w, v.Elem())
	}, nil
}

func newSliceEncoder(typ reflect.Type) (encoder, error) {
	elem := typ.Elem()
	for elem.Kind() == reflect.Pointer {
		elem = elem.Elem()
	}
	if elem.Kind() == reflect.Slice {
		return nil, fmt.Errorf("unsupported type: %s", typ)
	}

	enc, err := newFieldEncoder(typ.Elem())
	if err != nil {
		return nil, err
	}

	return func(w io.Writer, v reflect.Value) error {
		for i := range v.Len() {
			if i != 0 {
				w.Write(nl) //nolint:errcheck // It's a buffer so will always succeed.
			}
			if err := enc(w, v.Index(i)); err != nil {
				return err
			}
		}

		return nil
	}, nil
}

func newDateEncoder(reflect.Type) (encoder, error) {
	return func(w io.Writer, v reflect.Value) error {
		t := v.Interface().(time.Time) //nolint:forcetypeassert
		if t.IsZero() {
			return nil
		}
		_, err := w.Write(atob(strconv.Itoa(int(t.Unix()))))
		return err
	}, nil
}

func newIntEncoder(reflect.Type) (encoder, error) {
	return func(w io.Writer, v reflect.Value) error {
		i := v.Int()
		if i == 0 {
			return nil
		}
		_, err := w.Write(atob(strconv.FormatInt(i, 10)))
		return err
	}, nil
}

func newUintEncoder(reflect.Type) (encoder, error) {
	return func(w io.Writer, v reflect.Value) error {
		i := v.Uint()
		if i == 0 {
			return nil
		}
		_, err := w.Write(atob(strconv.FormatUint(i, 10)))
		return err
	}, nil
}

func newFloatEncoder(typ reflect.Type) (encoder, error) {
	bits := typ.Bits()
	return func(w io.Writer, v reflect.Value) error {
		f := v.Float()
		if f == 0 {
			return nil
		}
		_, err := w.Write(strconv.AppendFloat(nil, f, 'f', -1, bits))
		return err
	}, nil
}

func newByteArrayEncoder(typ reflect.Type) (encoder, error) {
	size := typ.Len()
	return func(w io.Writer, v reflect.Value) error {
		var isNonZero bool
		b := make([]byte, size)
		for i := range size {
			if n := v.Index(i).Uint(); n > 0 {
				b[i] = byte(n)
				isNonZero = true
			}
		}
		if !isNonZero {
			return nil
		}
		_, err := hex.NewEncoder(w).Write(b)
		return err
	}, nil
}

func newStringEncoder(reflect.Type) (encoder, error) {
	return func(w io.Writer, v reflect.Value) error {
		_, err := w.Write(atob(v.String()))
		return err
	}, nil
}
