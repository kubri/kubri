package desc_test

import (
	"errors"
	"io"
	"reflect"
	"strings"
	"testing"
	"testing/iotest"
	"time"

	"github.com/google/go-cmp/cmp"

	"github.com/kubri/kubri/integrations/arch/desc"
	"github.com/kubri/kubri/internal/test"
)

func TestUnmarshal(t *testing.T) {
	in := `%STRING%
test

%POINTER%
test

%SLICE%
foo
bar
baz

%HEX%
01020304

%INT%
1

%INT8%
1

%INT16%
1

%INT32%
1

%UINT%
1

%UINT8%
1

%UINT16%
1

%UINT32%
1

%FLOAT32%
1.123

%FLOAT64%
1.123

%MARSHALER%
test

%DATE%
1673377465`

	s := "test"

	want := record{
		String:    "test",
		Pointer:   &s,
		Slice:     []string{"foo", "bar", "baz"},
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
			msg: "pointer to struct pointer",
			in:  in,
			want: func() **record {
				r := &want
				return &r
			}(),
		},
		{
			msg: "multi-line text",
			in: `%STRING%
foo
bar
baz

foobar
%MARSHALER%
foo
bar
baz

foobar
`,
			want: record{
				String:    "foo\nbar\nbaz\n\nfoobar",
				Marshaler: &marshaler{"foo\nbar\nbaz\n\nfoobar"},
			},
		},
		{
			msg: "multi-line text starting with empty line",
			in: `%STRING%

foo
bar
%MARSHALER%

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
			in: `%STRING%
%HEX%
%INT%
%INT8%
%INT16%
%INT32%
%UINT%
%UINT8%
%UINT16%
%UINT32%
%FLOAT32%
%FLOAT64%
%MARSHALER%
%DATE%
`,
			want: record{},
		},
		{
			msg: "unknown keys",
			in: `%FOO%
bar
%STRING%
test
`,
			want: record{
				String: "test",
			},
		},
		{
			msg: "non-pointer unmarshaler",
			in: `%MARSHALER%
test
`,
			want: struct{ Marshaler marshaler }{
				Marshaler: marshaler{"test"},
			},
		},
		{
			msg: "unexported/ignored fields",
			in:  "%UNEXPORTED%\nfoo\n%IGNORED%\nbar\n%TEST%\nbaz\n",
			want: struct {
				unexported string
				Ignored    string `desc:"-"`
				Test       string
			}{Test: "baz"},
		},
		{
			msg: "named fields",
			in:  "%ALIAS%\nfoo\n",
			want: struct {
				Name string `desc:"ALIAS"`
			}{Name: "foo"},
		},
		{
			msg:  "leading blank lines",
			in:   "\n\n%STRING%\nfoo\n",
			want: record{String: "foo"},
		},
	}

	opts := test.ExportAll()

	for _, test := range tests {
		v := reflect.New(reflect.TypeOf(test.want))
		if err := desc.Unmarshal([]byte(test.in), v.Interface()); err != nil {
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
			msg:    "invalid date",
			reader: strings.NewReader("%DATE%\ntest\n"),
			value:  &record{},
			err:    `strconv.Atoi: parsing "test": invalid syntax`,
		},
		{
			msg:    "invalid integer",
			reader: strings.NewReader("%INT%\ntest\n"),
			value:  &record{},
			err:    `strconv.ParseInt: parsing "test": invalid syntax`,
		},
		{
			msg:    "invalid unsigned integer",
			reader: strings.NewReader("%UINT%\ntest\n"),
			value:  &record{},
			err:    `strconv.ParseUint: parsing "test": invalid syntax`,
		},
		{
			msg:    "invalid float",
			reader: strings.NewReader("%FLOAT64%\ntest\n"),
			value:  &record{},
			err:    `strconv.ParseFloat: parsing "test": invalid syntax`,
		},
		{
			msg:    "invalid hex data",
			reader: strings.NewReader("%HEX%\ntest\n"),
			value:  &record{},
			err:    "encoding/hex: invalid byte: U+0074 't'",
		},
		{
			msg:    "invalid hex length",
			reader: strings.NewReader("%HEX%\nFFFFFFFFFF\n"),
			value:  &record{},
			err:    "hex data would overflow byte array",
		},
		{
			msg:    "invalid value in slice",
			reader: strings.NewReader("%V%\ntest\n"),
			value:  &struct{ V []int }{},
			err:    `strconv.ParseInt: parsing "test": invalid syntax`,
		},
		{
			msg:    "missing field name",
			reader: strings.NewReader("Test\n"),
			value:  &record{},
			err:    `invalid line: "Test"`,
		},
		{
			msg:    "empty field name",
			reader: strings.NewReader("%%\n"),
			value:  &record{},
			err:    `invalid line: "%%"`,
		},
		{
			msg:    "unmarshaler error",
			reader: strings.NewReader("%MARSHALER%\ntest\n"),
			value: &struct{ Marshaler errMarshaler }{
				Marshaler: errMarshaler{errors.New("unmarshal error")},
			},
			err: "unmarshal error",
		},
		{
			msg:    "reader error on first read",
			reader: io.MultiReader(strings.NewReader("%STRING%\ntest"), iotest.ErrReader(errors.New("reader error"))),
			value:  &record{},
			err:    "reader error",
		},
		{
			msg:    "reader error on second read",
			reader: io.MultiReader(strings.NewReader("%STRING%\ntest\n"), iotest.ErrReader(errors.New("reader error"))),
			value:  &record{},
			err:    "reader error",
		},
	}

	opts := test.CompareErrorMessages()

	for _, test := range tests {
		err := desc.NewDecoder(test.reader).Decode(test.value)
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
		{"string", []byte("%V%\ntest"), &struct{ V string }{}},
		{"int", []byte("%V%\n1"), &struct{ V int }{}},
		{"uint", []byte("%V%\n1"), &struct{ V uint }{}},
		{"float64", []byte("%V%\n1"), &struct{ V float64 }{}},
		{"time.Time", []byte("%V%\n1673377465"), &struct{ V time.Time }{}},
		{"[8]byte", []byte("%V%\n0102030405060708"), &struct{ V [8]byte }{}},
	}

	for _, bm := range benchmarks {
		b.Run(bm.name, func(b *testing.B) {
			for b.Loop() {
				if err := desc.Unmarshal(bm.data, bm.v); err != nil {
					b.Fatal(err)
				}
			}
		})
	}
}
