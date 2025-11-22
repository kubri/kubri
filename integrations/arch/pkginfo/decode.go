package pkginfo

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

type (
	decoder      func(*bufio.Reader, reflect.Value) error
	fieldDecoder func([]byte, reflect.Value) error
)

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
		return dec.(decoder)(d.r, val) //nolint:forcetypeassert
	}

	dec, err := newDecoder(typ)
	if err != nil {
		return err
	}
	decoders.Store(typ, dec)
	return dec(d.r, val)
}

//nolint:gocognit
func newDecoder(typ reflect.Type) (decoder, error) {
	decoders := map[string]fieldDecoder{}

	var deref []reflect.Type
	for typ.Kind() == reflect.Pointer {
		typ = typ.Elem()
		deref = append(deref, typ)
	}

	if typ.Kind() != reflect.Struct {
		return nil, errors.New("must use pointer to struct")
	}

	for i := range typ.NumField() {
		field := typ.Field(i)

		name := getFieldName(field)
		if name == "" {
			continue
		}

		dec, err := newFieldDecoder(field.Type)
		if err != nil {
			return nil, err
		}

		decoders[name] = func(b []byte, v reflect.Value) error {
			return dec(b, v.Field(i))
		}
	}

	return func(r *bufio.Reader, v reflect.Value) error {
		for _, t := range deref {
			if v.IsNil() {
				v.Set(reflect.New(t))
			}
			v = v.Elem()
		}

		for {
			line, err := r.ReadSlice('\n')
			if err != nil && err != io.EOF {
				return err
			}
			if err == io.EOF && len(line) == 0 {
				return nil
			}
			if len(line) == 0 || line[0] == '#' || line[0] == '\n' {
				continue
			}
			key, val, ok := bytes.Cut(line, sep)
			if !ok {
				return fmt.Errorf("invalid line: %q", bytes.TrimRight(line, "\n"))
			}
			if dec := decoders[btoa(trim(key))]; dec != nil {
				if err := dec(trim(val), v); err != nil {
					return err
				}
			}
		}
	}, nil
}

//nolint:gochecknoglobals
var unmarshalerType = reflect.TypeFor[encoding.TextUnmarshaler]()

func newFieldDecoder(typ reflect.Type) (fieldDecoder, error) {
	switch {
	case typ == dateType:
		return newDateDecoder(typ)
	case typ.Implements(unmarshalerType), reflect.PointerTo(typ).Implements(unmarshalerType):
		return newUnmarshalerDecoder(typ)
	}

	switch typ.Kind() {
	case reflect.Ptr:
		return newPtrDecoder(typ)
	case reflect.Struct:
		return newStructDecoder(typ)
	case reflect.Slice:
		return newSliceDecoder(typ)
	case reflect.Bool:
		return newBoolDecoder(typ)
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

func newUnmarshalerDecoder(typ reflect.Type) (fieldDecoder, error) {
	mustAddr := reflect.PointerTo(typ).Implements(unmarshalerType)
	isPtr := typ.Kind() == reflect.Pointer
	if isPtr {
		typ = typ.Elem()
	}

	return func(b []byte, v reflect.Value) error {
		if len(b) == 0 {
			return nil
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

func newPtrDecoder(typ reflect.Type) (fieldDecoder, error) {
	typ = typ.Elem()
	dec, err := newFieldDecoder(typ)
	if err != nil {
		return nil, err
	}
	return func(b []byte, v reflect.Value) error {
		if v.IsNil() {
			v.Set(reflect.New(typ))
		}
		return dec(b, v.Elem())
	}, nil
}

func newSliceDecoder(typ reflect.Type) (fieldDecoder, error) {
	typ = typ.Elem()

	// Don't allow slice of structs or slices.
	fieldType := typ
	for fieldType.Kind() == reflect.Pointer {
		fieldType = fieldType.Elem()
	}
	if k := fieldType.Kind(); k == reflect.Struct || k == reflect.Slice {
		return nil, errors.New("deeply nested fields not supported")
	}

	dec, err := newFieldDecoder(typ)
	if err != nil {
		return nil, err
	}

	return func(b []byte, v reflect.Value) error {
		elem := reflect.New(typ).Elem()
		if err := dec(b, elem); err != nil {
			return err
		}
		v.Set(reflect.Append(v, elem))
		return nil
	}, nil
}

func newStructDecoder(typ reflect.Type) (fieldDecoder, error) {
	decoders := map[string]fieldDecoder{}

	for i := range typ.NumField() {
		field := typ.Field(i)

		// Don't allow nested structs or slices.
		fieldType := field.Type
		for fieldType.Kind() == reflect.Pointer {
			fieldType = fieldType.Elem()
		}
		if k := fieldType.Kind(); k == reflect.Struct || k == reflect.Slice {
			return nil, errors.New("deeply nested fields not supported")
		}

		name := getFieldName(field)
		if name == "" {
			continue
		}

		dec, err := newFieldDecoder(field.Type)
		if err != nil {
			return nil, err
		}

		decoders[name] = func(b []byte, v reflect.Value) error {
			return dec(b, v.Field(i))
		}
	}

	return func(b []byte, v reflect.Value) error {
		if len(b) == 0 {
			return nil
		}
		key, val, ok := bytes.Cut(b, equal)
		if !ok {
			return fmt.Errorf("invalid value: %q", b)
		}
		if dec := decoders[btoa(trim(key))]; dec != nil {
			if err := dec(trim(val), v); err != nil {
				return err
			}
		}
		return nil
	}, nil
}

func newDateDecoder(reflect.Type) (fieldDecoder, error) {
	return func(b []byte, v reflect.Value) error {
		if len(b) == 0 {
			return nil
		}
		n, err := strconv.Atoi(btoa(b))
		if err != nil {
			return err
		}
		t := time.Unix(int64(n), 0)
		if !t.IsZero() {
			v.Set(reflect.ValueOf(t))
		}
		return nil
	}, nil
}

func newBoolDecoder(reflect.Type) (fieldDecoder, error) {
	return func(b []byte, v reflect.Value) error {
		if len(b) == 0 {
			return nil
		}
		val, err := strconv.ParseBool(btoa(b))
		if err != nil {
			return err
		}
		v.SetBool(val)
		return nil
	}, nil
}

func newIntDecoder(typ reflect.Type) (fieldDecoder, error) {
	bits := typ.Bits()
	return func(b []byte, v reflect.Value) error {
		if len(b) == 0 {
			return nil
		}
		i, err := strconv.ParseInt(btoa(b), 10, bits)
		if err != nil {
			return err
		}
		v.SetInt(i)
		return nil
	}, nil
}

func newUintDecoder(typ reflect.Type) (fieldDecoder, error) {
	bits := typ.Bits()
	return func(b []byte, v reflect.Value) error {
		if len(b) == 0 {
			return nil
		}
		i, err := strconv.ParseUint(btoa(b), 10, bits)
		if err != nil {
			return err
		}
		v.SetUint(i)
		return nil
	}, nil
}

func newFloatDecoder(typ reflect.Type) (fieldDecoder, error) {
	bits := typ.Bits()
	return func(b []byte, v reflect.Value) error {
		if len(b) == 0 {
			return nil
		}
		i, err := strconv.ParseFloat(btoa(b), bits)
		if err != nil {
			return err
		}
		v.SetFloat(i)
		return nil
	}, nil
}

func newStringDecoder(reflect.Type) (fieldDecoder, error) {
	return func(b []byte, v reflect.Value) error {
		v.SetString(btoa(b))
		return nil
	}, nil
}

func newByteArrayDecoder(typ reflect.Type) (fieldDecoder, error) {
	size := typ.Len()
	return func(b []byte, v reflect.Value) error {
		if len(b) == 0 {
			return nil
		}
		if hex.DecodedLen(len(b)) > size {
			return errors.New("hex data would overflow byte array")
		}
		out := make([]byte, size)
		if _, err := hex.Decode(out, b); err != nil {
			return err
		}
		for i, x := range out {
			v.Index(i).SetUint(uint64(x))
		}
		return nil
	}, nil
}
