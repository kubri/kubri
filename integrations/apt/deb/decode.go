package deb

import (
	"bufio"
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

var decoders sync.Map //nolint:gochecknoglobals

type decoder func(*bufio.Reader, reflect.Value) error

func Unmarshal(b []byte, v any) error {
	return NewDecoder(bytes.NewReader(b)).Decode(v)
}

type Decoder struct{ r *bufio.Reader }

func NewDecoder(r io.Reader) *Decoder {
	return &Decoder{bufio.NewReader(r)}
}

func (d *Decoder) Decode(v any) error {
	if v == nil {
		return errors.New("unsupported type: nil")
	}

	val := reflect.ValueOf(v)
	typ := val.Type()

	if typ.Kind() != reflect.Pointer {
		return errors.New("must use pointer")
	}

	if dec, ok := decoders.Load(typ); ok {
		err := dec.(decoder)(d.r, val) //nolint:forcetypeassert
		if err == io.EOF {
			return errors.New("unexpected end of input")
		}
		return err
	}

	dec, err := newDecoder(typ)
	if err != nil {
		return err
	}
	decoders.Store(typ, dec)
	return dec(d.r, val)
}

//nolint:gochecknoglobals
var unmarshalerType = reflect.TypeFor[encoding.TextUnmarshaler]()

func newDecoder(typ reflect.Type) (decoder, error) {
	switch {
	case typ == dateType:
		return newDateDecoder(typ)
	case typ.Implements(unmarshalerType), reflect.PointerTo(typ).Implements(unmarshalerType):
		return newUnmarshalerDecoder(typ)
	}

	switch typ.Kind() {
	case reflect.Ptr:
		return newPtrDecoder(typ)
	case reflect.Slice:
		return newSliceDecoder(typ)
	case reflect.Struct:
		return newStructDecoder(typ)
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return newIntDecoder(typ)
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return newUintDecoder(typ)
	case reflect.Float32, reflect.Float64:
		return newFloatDecoder(typ)
	case reflect.String:
		return newStringDecoder(typ)
	case reflect.Array:
		if typ.Elem().Kind() == reflect.Uint8 {
			return newByteArrayDecoder(typ)
		}
	}

	return nil, fmt.Errorf("unsupported type: %s", typ)
}

func newUnmarshalerDecoder(typ reflect.Type) (decoder, error) {
	mustAddr := reflect.PointerTo(typ).Implements(unmarshalerType)
	isPtr := typ.Kind() == reflect.Pointer
	if isPtr {
		typ = typ.Elem()
	}

	return func(r *bufio.Reader, v reflect.Value) error {
		b, err := readlines(r)
		if err != nil || len(b) == 0 {
			return err
		}
		if isPtr && v.IsNil() {
			v.Set(reflect.New(typ))
		}
		if mustAddr {
			v = v.Addr()
		}
		return v.Interface().(encoding.TextUnmarshaler).UnmarshalText(b) //nolint:forcetypeassert
	}, nil
}

func newPtrDecoder(typ reflect.Type) (decoder, error) {
	typ = typ.Elem()
	dec, err := newDecoder(typ)
	if err != nil {
		return nil, err
	}
	return func(r *bufio.Reader, v reflect.Value) error {
		if v.IsNil() {
			v.Set(reflect.New(typ))
		}
		return dec(r, v.Elem())
	}, nil
}

func newSliceDecoder(typ reflect.Type) (decoder, error) {
	typ = typ.Elem()
	dec, err := newDecoder(typ)
	if err != nil {
		return nil, err
	}

	return func(r *bufio.Reader, v reflect.Value) error {
		v.SetLen(0)

		for {
			for {
				b, err := r.ReadByte()
				if err == io.EOF {
					return nil
				}
				if err != nil {
					return err
				}
				if b != '\n' && b != '\r' {
					_ = r.UnreadByte()
					break
				}
			}
			elem := reflect.New(typ).Elem()
			if err := dec(r, elem); err != nil {
				return err
			}
			v.Set(reflect.Append(v, elem))
		}
	}, nil
}

func newStructDecoder(typ reflect.Type) (decoder, error) {
	decoders := map[string]decoder{}

	for i := range typ.NumField() {
		field := typ.Field(i)

		name := getFieldName(field)
		if name == "" {
			continue
		}

		dec, err := newDecoder(field.Type)
		if err != nil {
			return nil, err
		}

		decoders[name] = func(r *bufio.Reader, v reflect.Value) error {
			return dec(r, v.Field(i))
		}
	}

	return func(r *bufio.Reader, v reflect.Value) error {
		for {
			c, err := r.Peek(1)
			if err != nil && err != io.EOF {
				return err
			}
			if err == io.EOF || c[0] == '\n' || c[0] == '\r' {
				return nil
			}
			b, err := r.ReadSlice(':')
			if err != nil {
				return err
			}
			if i := bytes.IndexByte(b, '\n'); i != -1 {
				return fmt.Errorf("invalid line: %q", b[:i])
			}
			if dec := decoders[btoa(trim(b)[:len(b)-1])]; dec != nil {
				if err := dec(r, v); err != nil {
					return err
				}
			} else {
				_, _ = readlines(r) // Discard value for unknown key.
			}
		}
	}, nil
}

func newDateDecoder(reflect.Type) (decoder, error) {
	return func(r *bufio.Reader, v reflect.Value) error {
		b, err := readline(r)
		if err != nil || len(b) == 0 {
			return err
		}
		t, err := time.Parse(time.RFC1123, btoa(b))
		if err != nil {
			return err
		}
		if !t.IsZero() {
			v.Set(reflect.ValueOf(t))
		}
		return nil
	}, nil
}

func newIntDecoder(typ reflect.Type) (decoder, error) {
	bits := typ.Bits()
	return func(r *bufio.Reader, v reflect.Value) error {
		b, err := readline(r)
		if err != nil || len(b) == 0 {
			return err
		}
		i, err := strconv.ParseInt(btoa(b), 10, bits)
		if err != nil {
			return err
		}
		v.SetInt(i)
		return nil
	}, nil
}

func newUintDecoder(typ reflect.Type) (decoder, error) {
	bits := typ.Bits()
	return func(r *bufio.Reader, v reflect.Value) error {
		b, err := readline(r)
		if err != nil || len(b) == 0 {
			return err
		}
		i, err := strconv.ParseUint(btoa(b), 10, bits)
		if err != nil {
			return err
		}
		v.SetUint(i)
		return nil
	}, nil
}

func newFloatDecoder(typ reflect.Type) (decoder, error) {
	bits := typ.Bits()
	return func(r *bufio.Reader, v reflect.Value) error {
		b, err := readline(r)
		if err != nil || len(b) == 0 {
			return err
		}
		i, err := strconv.ParseFloat(btoa(b), bits)
		if err != nil {
			return err
		}
		v.SetFloat(i)
		return nil
	}, nil
}

func newStringDecoder(reflect.Type) (decoder, error) {
	return func(r *bufio.Reader, v reflect.Value) error {
		b, err := readlines(r)
		if err != nil {
			return err
		}
		v.SetString(btoa(b))
		return nil
	}, nil
}

func newByteArrayDecoder(typ reflect.Type) (decoder, error) {
	size := typ.Len()
	return func(r *bufio.Reader, v reflect.Value) error {
		b, err := readline(r)
		if err != nil || len(b) == 0 {
			return err
		}
		if hex.DecodedLen(len(b)) > size {
			return errors.New("hex data would overflow byte array")
		}
		out := make([]byte, size)
		if _, err = hex.Decode(out, b); err != nil {
			return err
		}
		for i, x := range out {
			v.Index(i).SetUint(uint64(x))
		}
		return nil
	}, nil
}

func readline(r *bufio.Reader) ([]byte, error) {
	b, err := r.ReadSlice('\n')
	if err != nil {
		return nil, err
	}
	return trim(b), nil
}

func readlines(r *bufio.Reader) ([]byte, error) {
	b, err := readline(r)
	if err != nil {
		return nil, err
	}

	for {
		p, err := r.Peek(1)
		if err != nil && err != io.EOF {
			return nil, err
		}
		if err == io.EOF || p[0] != ' ' && p[0] != '\t' {
			break
		}
		l, err := readline(r)
		if err != nil {
			return nil, err
		}
		b = append(b, '\n')
		if len(l) != 1 || l[0] != '.' {
			b = append(b, l...)
		}
	}

	return b, nil
}
