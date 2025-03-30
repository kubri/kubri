package pkginfo_test

import (
	"errors"
	"io"
	"reflect"
	"strings"
	"testing"
	"testing/iotest"
	"time"

	"github.com/google/go-cmp/cmp"

	"github.com/kubri/kubri/integrations/arch/pkginfo"
	"github.com/kubri/kubri/internal/test"
)

func TestUnmarshal(t *testing.T) {
	in := `# This is a comment
string = test
pointer = test
struct = string=test
struct = int=1
struct = bool=true
struct = unexported=foo
struct = ignored=bar
slice = foo
slice = bar
slice = baz
hex = 01020304
bool = true
int = 1
int8 = 1
int16 = 1
int32 = 1
uint = 1
uint8 = 1
uint16 = 1
uint32 = 1
float32 = 1.123
float64 = 1.123
marshaler = test
date = 1743356936
`

	s := "test"

	want := record{
		String:  "test",
		Pointer: &s,
		Struct: recordStruct{
			String: "test",
			Int:    1,
			Bool:   true,
		},
		Slice:     []string{"foo", "bar", "baz"},
		Hex:       [4]byte{1, 2, 3, 4},
		Bool:      true,
		Int:       1,
		Int8:      1,
		Int16:     1,
		Int32:     1,
		Uint:      1,
		Uint8:     1,
		Uint16:    1,
		Uint32:    1,
		Float32:   1.123,
		Float64:   1.123,
		Marshaler: &marshaler{"test"},
		Date:      time.Date(2025, 3, 30, 17, 48, 56, 0, time.UTC),
	}

	tests := []struct {
		msg  string
		in   string
		want any
	}{
		{
			msg:  "struct",
			in:   in,
			want: want,
		},
		{
			msg:  "struct pointer",
			in:   in,
			want: &want,
		},
		{
			msg: "pointer to struct pointer",
			in:  in,
			want: func() **record {
				r := &want
				return &r
			}(),
		},
		{
			msg: "empty values",
			in: `string = 
struct = 
hex = 
bool = 
int = 
int8 = 
int16 = 
int32 = 
uint = 
uint8 = 
uint16 = 
uint32 = 
float32 = 
float64 = 
marshaler = 
date = 
`,
			want: record{},
		},
		{
			msg: "unknown keys",
			in:  "foo = bar\nstring = test\n",
			want: record{
				String: "test",
			},
		},
		{
			msg: "unexported/ignored fields",
			in:  "unexported = foo\nignored = bar\ntest = baz\n",
			want: struct {
				unexported string
				Ignored    string `pkginfo:"-"`
				Test       string
			}{Test: "baz"},
		},
		{
			msg: "named fields",
			in:  "alias = foo\n",
			want: struct {
				Name string `pkginfo:"alias"`
			}{Name: "foo"},
		},
		{
			msg: "funky spacing",
			in:  "\nstring  =  test\n\n	int = 1\n",
			want: record{
				String: "test",
				Int:    1,
			},
		},
	}

	opts := test.ExportAll()

	for _, test := range tests {
		v := reflect.New(reflect.TypeOf(test.want))
		if err := pkginfo.Unmarshal([]byte(test.in), v.Interface()); err != nil {
			t.Error(test.msg, err)
		} else {
			if diff := cmp.Diff(test.want, v.Elem().Interface(), opts); diff != "" {
				t.Errorf("%s:\n%s", test.msg, diff)
			}
		}
	}
}

func TestDecodeErrors(t *testing.T) {
	tests := []struct {
		msg    string
		reader io.Reader
		value  any
		err    string
	}{
		{
			msg:   "non-pointer",
			value: record{},
			err:   "must use pointer",
		},
		{
			msg:   "non-struct-pointer",
			value: (*string)(nil),
			err:   "must use pointer to struct",
		},
		{
			msg:   "deeply nested struct",
			value: &struct{ V struct{ V *struct{} } }{},
			err:   "deeply nested fields not supported",
		},
		{
			msg:   "deeply nested slice",
			value: &struct{ V struct{ V []string } }{},
			err:   "deeply nested fields not supported",
		},
		{
			msg:   "slice of structs",
			value: &struct{ V []*struct{} }{},
			err:   "deeply nested fields not supported",
		},
		{
			msg:    "struct reader error",
			reader: iotest.ErrReader(errors.New("reader error")),
			value:  &record{},
			err:    "reader error",
		},
		{
			msg:   "nil",
			value: nil,
			err:   "unsupported type: nil",
		},
		{
			msg:   "unsupported type",
			value: &struct{ V complex128 }{},
			err:   "unsupported type: complex128",
		},
		{
			msg:   "unsupported pointer type",
			value: &struct{ V *complex128 }{},
			err:   "unsupported type: complex128",
		},
		{
			msg:   "unsupported type in slice",
			value: &struct{ V []complex128 }{},
			err:   "unsupported type: complex128",
		},
		{
			msg:   "unsupported type in struct",
			value: &struct{ V struct{ V complex128 } }{},
			err:   "unsupported type: complex128",
		},
		{
			msg:    "invalid date",
			reader: strings.NewReader("date = test\n"),
			value:  &record{},
			err:    `strconv.Atoi: parsing "test": invalid syntax`,
		},
		{
			msg:    "invalid bool",
			reader: strings.NewReader("bool = test\n"),
			value:  &record{},
			err:    `strconv.ParseBool: parsing "test": invalid syntax`,
		},
		{
			msg:    "invalid integer",
			reader: strings.NewReader("int = test\n"),
			value:  &record{},
			err:    `strconv.ParseInt: parsing "test": invalid syntax`,
		},
		{
			msg:    "invalid unsigned integer",
			reader: strings.NewReader("uint = test\n"),
			value:  &record{},
			err:    `strconv.ParseUint: parsing "test": invalid syntax`,
		},
		{
			msg:    "invalid float",
			reader: strings.NewReader("float64 = test\n"),
			value:  &record{},
			err:    `strconv.ParseFloat: parsing "test": invalid syntax`,
		},
		{
			msg:    "invalid hex data",
			reader: strings.NewReader("hex = test\n"),
			value:  &record{},
			err:    "encoding/hex: invalid byte: U+0074 't'",
		},
		{
			msg:    "invalid hex length",
			reader: strings.NewReader("hex = FFFFFFFFFF\n"),
			value:  &record{},
			err:    "hex data would overflow byte array",
		},
		{
			msg:    "invalid struct data",
			reader: strings.NewReader("struct = test\n"),
			value:  &record{},
			err:    `invalid value: "test"`,
		},
		{
			msg:    "invalid value in slice",
			reader: strings.NewReader("v = test\n"),
			value:  &struct{ V []int }{},
			err:    `strconv.ParseInt: parsing "test": invalid syntax`,
		},
		{
			msg:    "invalid value in struct",
			reader: strings.NewReader("v = v=test\n"),
			value:  &struct{ V struct{ V int } }{},
			err:    `strconv.ParseInt: parsing "test": invalid syntax`,
		},
		{
			msg:    "invalid line",
			reader: strings.NewReader("foo"),
			value:  &record{},
			err:    `invalid line: "foo"`,
		},
		{
			msg:    "unmarshaler error",
			reader: strings.NewReader("marshaler = test\n"),
			value: &struct{ Marshaler errMarshaler }{
				Marshaler: errMarshaler{errors.New("unmarshal error")},
			},
			err: "unmarshal error",
		},
	}

	opts := test.CompareErrorMessages()

	for _, test := range tests {
		err := pkginfo.NewDecoder(test.reader).Decode(test.value)
		if diff := cmp.Diff(errors.New(test.err), err, opts); diff != "" {
			t.Errorf("%s returned unexpected error:\n%s", test.msg, diff)
		}
	}
}

func BenchmarkUnmarshal(b *testing.B) {
	benchmarks := []struct {
		name string
		data []byte
		v    any
	}{
		{"string", []byte("v = test"), &struct{ V string }{}},
		{"int", []byte("v = 1"), &struct{ V int }{}},
		{"uint", []byte("v = 1"), &struct{ V uint }{}},
		{"float64", []byte("v = 1"), &struct{ V float64 }{}},
		{"time.Time", []byte("v = 1673377465"), &struct{ V time.Time }{}},
		{"[8]byte", []byte("v = 00000000499602D2"), &struct{ V [8]byte }{}},
	}

	for _, bm := range benchmarks {
		b.Run(bm.name, func(b *testing.B) {
			for b.Loop() {
				if err := pkginfo.Unmarshal(bm.data, bm.v); err != nil {
					b.Fatal(err)
				}
			}
		})
	}
}
