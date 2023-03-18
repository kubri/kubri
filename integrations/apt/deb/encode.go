package deb

import (
	"bytes"
	"encoding"
	"encoding/hex"
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

//nolint:gochecknoglobals
var (
	stringerType  = reflect.TypeOf((*fmt.Stringer)(nil)).Elem()
	marshalerType = reflect.TypeOf((*encoding.TextMarshaler)(nil)).Elem()
)

func newEncoder(typ reflect.Type) (encoder, error) {
	if typ == dateType {
		return newDateEncoder(typ)
	}

	switch {
	case typ.Implements(stringerType):
		return newStringerEncoder(typ)
	case typ.Implements(marshalerType):
		return newMarshalerEncoder(typ)
	}

	switch typ.Kind() {
	case reflect.Ptr:
		return newPtrEncoder(typ)
	case reflect.Slice:
		return newSliceEncoder(typ)
	case reflect.Struct:
		return newStructEncoder(typ)
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

func newStringerEncoder(reflect.Type) (encoder, error) {
	return func(w io.Writer, v reflect.Value) error {
		s := v.Interface().(fmt.Stringer).String() //nolint:forcetypeassert
		if len(s) > 0 {
			_, err := w.Write([]byte(s))
			return err
		}
		return nil
	}, nil
}

func newMarshalerEncoder(typ reflect.Type) (encoder, error) {
	ptr := typ.Kind() == reflect.Pointer
	return func(w io.Writer, v reflect.Value) error {
		if ptr && v.IsNil() {
			return nil
		}
		b, err := v.Interface().(encoding.TextMarshaler).MarshalText()
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
	enc, err := newEncoder(typ)
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
	enc, err := newEncoder(typ.Elem())
	if err != nil {
		return nil, err
	}

	return func(w io.Writer, v reflect.Value) error {
		for i := 0; i < v.Len(); i++ {
			if i != 0 {
				if _, err := w.Write(nl); err != nil {
					return err
				}
			}
			if err := enc(w, v.Index(i)); err != nil {
				return err
			}
		}

		return nil
	}, nil
}

//nolint:gocognit
func newStructEncoder(typ reflect.Type) (encoder, error) {
	encoders := []func(buf *bytes.Buffer, w io.Writer, v reflect.Value) error{}

	for i := 0; i < typ.NumField(); i++ {
		i := i
		field := typ.Field(i)

		n := getFieldName(field)
		if n == "" {
			continue
		}
		name := []byte(n)

		enc, err := newEncoder(field.Type)
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
			if _, err = w.Write(name); err != nil {
				return err
			}
			if _, err = w.Write(colon); err != nil {
				return err
			}

			// Only write colon after space if content doesn't start with line break.
			if c, _ := buf.ReadByte(); c != '\n' {
				if _, err = w.Write(space); err != nil {
					return err
				}
			}
			_ = buf.UnreadByte()

			if _, err = io.Copy(w, buf); err != nil {
				return err
			}
			_, err = w.Write(nl)
			return err
		})
	}

	return func(w io.Writer, v reflect.Value) error {
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

func newDateEncoder(reflect.Type) (encoder, error) {
	return func(w io.Writer, v reflect.Value) error {
		t := v.Interface().(time.Time) //nolint:forcetypeassert
		if t.IsZero() {
			return nil
		}
		_, err := w.Write([]byte(t.Format(time.RFC1123)))
		return err
	}, nil
}

func newIntEncoder(reflect.Type) (encoder, error) {
	return func(w io.Writer, v reflect.Value) error {
		i := v.Int()
		if i == 0 {
			return nil
		}
		_, err := w.Write([]byte(strconv.FormatInt(i, 10)))
		return err
	}, nil
}

func newUintEncoder(reflect.Type) (encoder, error) {
	return func(w io.Writer, v reflect.Value) error {
		i := v.Uint()
		if i == 0 {
			return nil
		}
		_, err := w.Write([]byte(strconv.FormatUint(i, 10)))
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
		_, err := w.Write([]byte(strconv.FormatFloat(f, 'f', -1, bits)))
		return err
	}, nil
}

func newByteArrayEncoder(typ reflect.Type) (encoder, error) {
	size := typ.Len()
	return func(w io.Writer, v reflect.Value) error {
		var isNonZero bool
		b := make([]byte, size)
		for i := 0; i < size; i++ {
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
		in := []byte(v.String())

		if i := bytes.IndexByte(in, '\n'); i == -1 {
			_, err := w.Write(in)
			return err
		}

		out := make([]byte, 0, len(in))
		for i, c := range in {
			if c == '\n' {
				var last byte
				for j := i - 1; j > 0; j-- {
					if l := in[j]; l != '\n' && l != '\r' {
						break
					}
					last = in[j]
				}
				if last == '\n' {
					out = append(out, '.')
				}
				out = append(out, '\n', ' ')
			} else {
				out = append(out, c)
			}
		}

		_, err := w.Write(out)
		return err
	}, nil
}
