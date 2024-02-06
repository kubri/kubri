package deb_test

import (
	"errors"
	"io"
	"reflect"
	"strings"
	"testing"
	"testing/iotest"
	"time"

	"github.com/google/go-cmp/cmp"

	"github.com/kubri/kubri/integrations/apt/deb"
	"github.com/kubri/kubri/internal/test"
)

func TestUnmarshal(t *testing.T) {
	in := `String: test
Hex: 01020304
Int: 1
Int8: 1
Int16: 1
Int32: 1
Uint: 1
Uint8: 1
Uint16: 1
Uint32: 1
Float32: 1.123
Float64: 1.123
Marshaler: test
Date: Tue, 10 Jan 2023 19:04:25 UTC
`

	want := record{
		String:    "test",
		Hex:       [4]byte{1, 2, 3, 4},
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
		Date:      time.Date(2023, 1, 10, 19, 4, 25, 0, time.UTC),
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
			msg:  "struct slice",
			in:   in + "\r\n\r\n" + in,
			want: []record{want, want},
		},
		{
			msg:  "struct pointer slice",
			in:   in + "\n" + in,
			want: []*record{&want, &want},
		},
		{
			msg: "multi-line text",
			in: `String: foo
 bar
 baz
 .
 foobar
Marshaler: foo
 bar
 baz
 .
 foobar
`,
			want: record{
				String:    "foo\nbar\nbaz\n\nfoobar",
				Marshaler: &marshaler{"foo\nbar\nbaz\n\nfoobar"},
			},
		},
		{
			msg: "multi-line text starting with empty line",
			in: `String:
 foo
 bar
Marshaler:
 foo
 bar
`,
			want: record{
				String:    "\nfoo\nbar",
				Marshaler: &marshaler{"\nfoo\nbar"},
			},
		},
		{
			msg: "empty values",
			in: `String: 
Hex: 
Int: 
Int8: 
Int16: 
Int32: 
Uint: 
Uint8: 
Uint16: 
Uint32: 
Float32: 
Float64: 
Marshaler: 
Date: 
`,
			want: record{},
		},
		{
			msg: "unknown keys",
			in: `Foo: bar
String: test
`,
			want: record{
				String: "test",
			},
		},
		{
			msg: "non-pointer unmarshaler",
			in: `Marshaler: test
`,
			want: struct{ Marshaler marshaler }{
				Marshaler: marshaler{"test"},
			},
		},
		{
			msg: "unexported/ignored fields",
			in:  "unexported: foo\nIgnored: bar\nTest: baz\n",
			want: struct {
				unexported string
				Ignored    string `deb:"-"`
				Test       string
			}{Test: "baz"},
		},
		{
			msg: "named fields",
			in:  "Alias: foo\n",
			want: struct {
				Name string `deb:"Alias"` //nolint:tagliatelle
			}{Name: "foo"},
		},
	}

	opts := test.ExportAll()

	for _, test := range tests {
		v := reflect.New(reflect.TypeOf(test.want))
		if err := deb.Unmarshal([]byte(test.in), v.Interface()); err != nil {
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
			msg:    "struct reader error",
			reader: iotest.ErrReader(errors.New("reader error")),
			value:  &record{},
			err:    "reader error",
		},
		{
			msg:    "slice reader error",
			reader: iotest.ErrReader(errors.New("reader error")),
			value:  &[]record{{}},
			err:    "reader error",
		},
		{
			msg:   "nil",
			value: nil,
			err:   "unsupported type: nil",
		},
		{
			msg:   "unsupported type",
			value: &[]struct{ V complex128 }{},
			err:   "unsupported type: complex128",
		},
		{
			msg:    "invalid date",
			reader: strings.NewReader("Date: test\n"),
			value:  &[]record{},
			err:    `parsing time "test" as "Mon, 02 Jan 2006 15:04:05 MST": cannot parse "test" as "Mon"`,
		},
		{
			msg:    "invalid integer",
			reader: strings.NewReader("Int: test\n"),
			value:  &[]record{},
			err:    `strconv.ParseInt: parsing "test": invalid syntax`,
		},
		{
			msg:    "invalid unsigned integer",
			reader: strings.NewReader("Uint: test\n"),
			value:  &[]record{},
			err:    `strconv.ParseUint: parsing "test": invalid syntax`,
		},
		{
			msg:    "invalid float",
			reader: strings.NewReader("Float64: test\n"),
			value:  &[]record{},
			err:    `strconv.ParseFloat: parsing "test": invalid syntax`,
		},
		{
			msg:    "invalid hex data",
			reader: strings.NewReader("Hex: test\n"),
			value:  &[]record{},
			err:    "encoding/hex: invalid byte: U+0074 't'",
		},
		{
			msg:    "invalid hex length",
			reader: strings.NewReader("Hex: FFFFFFFFFF\n"),
			value:  &[]record{},
			err:    "hex data would overflow byte array",
		},
		{
			msg:    "invalid line",
			reader: strings.NewReader("Foo\nBar:"),
			value:  &[]record{},
			err:    `invalid line: "Foo"`,
		},
		{
			msg:    "missing colon",
			reader: strings.NewReader("String"),
			value:  &[]record{},
			err:    "unexpected end of input",
		},
		{
			msg:    "missing line feed",
			reader: strings.NewReader("String:"),
			value:  &[]record{},
			err:    "unexpected end of input",
		},
		{
			msg:    "unmarshaler error",
			reader: strings.NewReader("Marshaler: test\n"),
			value: &struct{ Marshaler errMarshaler }{
				Marshaler: errMarshaler{errors.New("unmarshal error")},
			},
			err: "unmarshal error",
		},
	}

	opts := test.CompareErrorMessages()

	for _, test := range tests {
		err := deb.NewDecoder(test.reader).Decode(test.value)
		if diff := cmp.Diff(errors.New(test.err), err, opts); diff != "" {
			t.Errorf("%s returned unexpected error:\n%s", test.msg, diff)
		}
	}
}

func BenchmarkUnmarshal(b *testing.B) {
	type record struct {
		String string
		Hex    [4]byte
		Int    int
	}

	in := []byte(`String: foo
 bar
 baz
Hex: 01020304
Int: 1

String: foo
 bar
 baz
Hex: 01020304
Int: 1
`)

	var v []record

	for i := 0; i < b.N; i++ {
		deb.Unmarshal(in, &v)
	}
}
