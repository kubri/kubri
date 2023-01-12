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

type decoder func(io.Reader, reflect.Value) error

func Unmarshal(b []byte, v any) error {
	return NewDecoder(bytes.NewReader(b)).Decode(v)
}

type Decoder struct{ r io.Reader }

func NewDecoder(r io.Reader) *Decoder {
	return &Decoder{r}
}

func (d *Decoder) Decode(v any) error {
	val := reflect.ValueOf(v)
	typ := val.Type()

	if typ.Kind() != reflect.Pointer {
		return errors.New("must use pointer")
	}

	if dec, ok := decoders.Load(typ); ok {
		if err := dec.(decoder)(d.r, val); err != nil && err != io.EOF { //nolint:forcetypeassert
			return err
		}
		return nil
	}

	dec, err := newDecoder(typ)
	if err != nil {
		return err
	}
	decoders.Store(typ, dec)
	if err = dec(d.r, val); err != nil && err != io.EOF {
		return err
	}
	return nil
}

//nolint:gochecknoglobals
var unmarshalerType = reflect.TypeOf((*encoding.TextUnmarshaler)(nil)).Elem()

func newDecoder(typ reflect.Type) (decoder, error) {
	if typ == dateType {
		return newDateDecoder(typ)
	}

	if typ.Implements(unmarshalerType) {
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
	isPtr := typ.Kind() == reflect.Pointer

	return func(r io.Reader, v reflect.Value) error {
		b, err := io.ReadAll(r)
		if err != nil {
			return err
		}

		if len(b) == 0 {
			return nil
		}

		if isPtr && v.IsNil() {
			v.Set(reflect.New(typ.Elem()))
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
	return func(r io.Reader, v reflect.Value) error {
		if v.IsNil() {
			v.Set(reflect.New(typ))
		}
		return dec(r, v.Elem())
	}, nil
}

//nolint:gocognit,funlen
func newSliceDecoder(typ reflect.Type) (decoder, error) {
	typ = typ.Elem()
	dec, err := newDecoder(typ)
	if err != nil {
		return nil, err
	}

	return func(r io.Reader, v reflect.Value) error {
		v.SetLen(0)

		buf := bufPool.Get().(*bytes.Buffer) //nolint:forcetypeassert
		buf.Reset()
		defer bufPool.Put(buf)

		next := make([]byte, 1)
		br := false

		parse := func() error {
			elem := reflect.New(typ).Elem()
			if err = dec(buf, elem); err != nil && err != io.EOF {
				return err
			}
			v.Set(reflect.Append(v, elem))
			buf.Reset()
			return nil
		}

		for {
			_, err := r.Read(next)
			if err == io.EOF && buf.Len() > 0 {
				if err = parse(); err != nil {
					return err
				}
				return io.EOF
			}
			if err != nil {
				return err
			}

		SWITCH:
			switch next[0] {
			case '\r':
				continue
			case '\n':
				if br { //nolint:nestif
					if err = parse(); err != nil {
						return err
					}
					// Keep reading until we get the next item.
					for {
						_, err := r.Read(next)
						if err != nil {
							return err
						}
						if c := next[0]; c != '\n' && c != '\r' {
							goto SWITCH
						}
					}
				} else {
					br = true
					if _, err := buf.Write(next); err != nil {
						return err
					}
				}
			default:
				br = false
				if _, err := buf.Write(next); err != nil {
					return err
				}
			}
		}
	}, nil
}

func newStructDecoder(typ reflect.Type) (decoder, error) {
	decoders := map[string]decoder{}

	for i := 0; i < typ.NumField(); i++ {
		i := i
		field := typ.Field(i)

		name := getFieldName(field)
		if name == "" {
			continue
		}

		dec, err := newDecoder(field.Type)
		if err != nil {
			return nil, err
		}

		decoders[name] = func(r io.Reader, v reflect.Value) error {
			return dec(r, v.Field(i))
		}
	}

	return func(r io.Reader, v reflect.Value) error {
		rl := unmarshaler{bufio.NewReader(r)}

		buf := bufPool.Get().(*bytes.Buffer) //nolint:forcetypeassert
		buf.Reset()
		defer bufPool.Put(buf)

		for {
			line, err := rl.ReadLine()
			if err != nil {
				return err
			}
			if len(line) == 0 || bytes.Equal(line, nl) {
				return io.EOF
			}
			key, val, ok := bytes.Cut(line, colon)
			if !ok {
				return errors.New("invalid line: " + string(line))
			}

			if dec := decoders[string(trim(key))]; dec != nil {
				if _, err = buf.Write(trim(val)); err != nil {
					return err
				}
				if err = dec(buf, v); err != nil {
					return err
				}
				buf.Reset()
			}
		}
	}, nil
}

func newDateDecoder(typ reflect.Type) (decoder, error) {
	return func(r io.Reader, v reflect.Value) error {
		b, err := io.ReadAll(r)
		if err != nil {
			return err
		}
		t, err := time.Parse(time.RFC1123, string(b))
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
	return func(r io.Reader, v reflect.Value) error {
		b, err := io.ReadAll(r)
		if err != nil {
			return err
		}
		i, err := strconv.ParseInt(string(b), 10, bits)
		if err != nil {
			return err
		}
		v.SetInt(i)
		return nil
	}, nil
}

func newUintDecoder(typ reflect.Type) (decoder, error) {
	bits := typ.Bits()
	return func(r io.Reader, v reflect.Value) error {
		b, err := io.ReadAll(r)
		if err != nil {
			return err
		}
		i, err := strconv.ParseUint(string(b), 10, bits)
		if err != nil {
			return err
		}
		v.SetUint(i)
		return nil
	}, nil
}

func newFloatDecoder(typ reflect.Type) (decoder, error) {
	bits := typ.Bits()
	return func(r io.Reader, v reflect.Value) error {
		b, err := io.ReadAll(r)
		if err != nil {
			return err
		}
		i, err := strconv.ParseFloat(string(b), bits)
		if err != nil {
			return err
		}
		v.SetFloat(i)
		return nil
	}, nil
}

func newStringDecoder(typ reflect.Type) (decoder, error) {
	return func(r io.Reader, v reflect.Value) error {
		b, err := io.ReadAll(r)
		if err != nil {
			return err
		}
		v.SetString(string(b))
		return nil
	}, nil
}

func newByteArrayDecoder(typ reflect.Type) (decoder, error) {
	size := typ.Len()
	return func(r io.Reader, v reflect.Value) error {
		b, err := io.ReadAll(r)
		if err != nil {
			return err
		}

		if _, err = hex.Decode(b, b); err != nil {
			return err
		}
		for i := 0; i < size; i++ {
			v.Index(i).SetUint(uint64(b[i]))
		}
		return nil
	}, nil
}

type unmarshaler struct{ r *bufio.Reader }

func (u *unmarshaler) ReadLine() ([]byte, error) {
	b, err := u.r.ReadBytes('\n')
	if err != nil {
		return nil, err
	}

	for {
		p, err := u.r.Peek(1)
		if err != nil || p[0] != ' ' && p[0] != '\t' {
			break
		}

		l, err := u.r.ReadBytes('\n')
		if err != nil {
			return nil, err
		}

		if l = trim(l); len(l) == 0 {
			break
		}

		if bytes.Equal(dot, l) {
			b = append(b, '\n')
		} else {
			b = append(b, l...)
			b = append(b, '\n')
		}
	}

	return b, nil //nolint:nilerr
}
