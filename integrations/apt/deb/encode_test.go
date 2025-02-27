package deb_test

import (
	"errors"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"

	"github.com/kubri/kubri/integrations/apt/deb"
	"github.com/kubri/kubri/internal/test"
)

func TestMarshal(t *testing.T) {
	in := record{
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
		Stringer:  stringer{"test"},
		Marshaler: &marshaler{"test"},
		Date:      time.Date(2023, 1, 10, 19, 4, 25, 0, time.UTC),
	}

	want := `String: test
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
Stringer: test
Marshaler: test
Date: Tue, 10 Jan 2023 19:04:25 UTC
`

	tests := []struct {
		msg  string
		in   any
		want string
	}{
		{
			msg:  "struct",
			in:   in,
			want: want,
		},
		{
			msg:  "struct pointer",
			in:   &in,
			want: want,
		},
		{
			msg:  "struct slice",
			in:   []record{in, in},
			want: want + "\n" + want,
		},
		{
			msg:  "struct pointer slice",
			in:   []*record{&in, &in},
			want: want + "\n" + want,
		},
		{
			msg: "multi-line text",
			in: record{
				String:    "foo\nbar\nbaz\n\nfoobar",
				Stringer:  stringer{"foo\nbar\nbaz\n\nfoobar"},
				Marshaler: &marshaler{"foo\nbar\nbaz\n\nfoobar"},
			},
			want: `String: foo
 bar
 baz
 .
 foobar
Stringer: foo
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
		},
		{
			msg: "multi-line text starting with empty line",
			in: record{
				String:    "\nfoo\nbar",
				Stringer:  stringer{"\nfoo\nbar"},
				Marshaler: &marshaler{"\nfoo\nbar"},
			},
			want: `String:
 foo
 bar
Stringer:
 foo
 bar
Marshaler:
 foo
 bar
`,
		},
		{
			msg: "nil values",
			in: struct {
				Date      *time.Time
				Marshaler *marshaler
				Stringer  *stringer
				String    *string
				Int       *int
				Uint      *uint
				Float     *float64
			}{},
			want: "",
		},
		{
			msg: "zero values",
			in: struct {
				Marshaler *marshaler
				Stringer  *stringer
				String    *string
			}{
				Marshaler: &marshaler{},
				Stringer:  &stringer{},
				String:    new(string),
			},
			want: "",
		},
		{
			msg: "unexported/ignored fields",
			in: struct {
				unexported string
				Ignored    string `deb:"-"`
				Test       string
			}{unexported: "foo", Ignored: "bar", Test: "baz"},
			want: "Test: baz\n",
		},
		{
			msg: "named fields",
			in: struct {
				Name string `deb:"Alias"` //nolint:tagliatelle
			}{Name: "foo"},
			want: "Alias: foo\n",
		},
	}

	for _, test := range tests {
		b, err := deb.Marshal(test.in)
		if err != nil {
			t.Error(test.msg, err)
		}

		if diff := cmp.Diff(test.want, string(b)); diff != "" {
			t.Errorf("%s:\n%s", test.msg, diff)
		}
	}
}

func TestEncodeErrors(t *testing.T) {
	tests := []struct {
		msg   string
		value any
		err   string
	}{
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
			msg:   "marshaler error",
			value: &[]struct{ V *errMarshaler }{{V: &errMarshaler{errors.New("marshal error")}}},
			err:   "marshal error",
		},
	}

	opts := test.CompareErrorMessages()

	for _, test := range tests {
		_, err := deb.Marshal(test.value)
		if diff := cmp.Diff(errors.New(test.err), err, opts); diff != "" {
			t.Errorf("%s returned unexpected error:\n%s", test.msg, diff)
		}
	}

	// Ensure write errors are passed though from each part of the application.
	// Test up to 6 writes: field name, colon, space, value, newline (field end), newline (slice element end)
	for i := 1; i <= 6; i++ {
		want := errors.New("custom error")
		w := &errWriter{i, want}
		err := deb.NewEncoder(w).Encode([]record{{String: "foo"}, {String: "bar"}})
		if diff := cmp.Diff(err, want, cmpopts.EquateErrors()); diff != "" {
			t.Errorf("write %d should return error:\n%s", i, diff)
		}
	}
}

func BenchmarkMarshal(b *testing.B) {
	type record struct {
		String string
		Hex    [4]byte
		Int    int
	}

	v := []record{
		{
			String: "foo\nbar\nbaz",
			Hex:    [4]byte{1, 2, 3, 4},
			Int:    1,
		},
		{
			String: "foo\nbar\nbaz",
			Hex:    [4]byte{1, 2, 3, 4},
			Int:    1,
		},
	}

	for range b.N {
		deb.Marshal(v)
	}
}

type errWriter struct {
	n int
	e error
}

func (w *errWriter) Write(p []byte) (int, error) {
	w.n--
	if w.n == 0 {
		return 0, w.e
	}
	return len(p), nil
}
