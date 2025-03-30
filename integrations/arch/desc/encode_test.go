package desc_test

import (
	"errors"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"

	"github.com/kubri/kubri/integrations/arch/desc"
	"github.com/kubri/kubri/internal/test"
)

func TestMarshal(t *testing.T) {
	s := "test"

	in := record{
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
		Stringer:  &stringer{"test"},
		Marshaler: &marshaler{"test"},
		Date:      time.Date(2023, 1, 10, 19, 4, 25, 0, time.UTC),
	}

	want := `%STRING%
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

%STRINGER%
test

%MARSHALER%
test

%DATE%
1673377465

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
			msg:  "empty struct pointer",
			in:   &record{},
			want: "",
		},
		{
			msg:  "nil struct pointer",
			in:   (*record)(nil),
			want: "",
		},
		{
			msg: "zero values",
			in: record{
				Marshaler: &marshaler{},
				Stringer:  &stringer{},
				Pointer:   new(string),
			},
			want: "",
		},
		{
			msg: "unexported/ignored fields",
			in: struct {
				unexported string
				Ignored    string `desc:"-"`
				Test       string
			}{unexported: "foo", Ignored: "bar", Test: "baz"},
			want: "%TEST%\nbaz\n\n",
		},
		{
			msg: "named fields",
			in: struct {
				Name string `desc:"ALIAS"`
			}{Name: "foo"},
			want: "%ALIAS%\nfoo\n\n",
		},
	}

	for _, test := range tests {
		b, err := desc.Marshal(test.in)
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
			value: "",
			err:   "unsupported type: string",
		},
		{
			msg:   "unsupported field type",
			value: struct{ V complex128 }{},
			err:   "unsupported type: complex128",
		},
		{
			msg:   "unsupported pointer type",
			value: struct{ V *complex128 }{},
			err:   "unsupported type: complex128",
		},
		{
			msg:   "unsupported slice type",
			value: struct{ V []complex128 }{},
			err:   "unsupported type: complex128",
		},
		{
			msg:   "slice of slice",
			value: struct{ V [][]string }{},
			err:   "unsupported type: [][]string",
		},
		{
			msg:   "marshaler error",
			value: struct{ V *errMarshaler }{V: &errMarshaler{errors.New("marshal error")}},
			err:   "marshal error",
		},
		{
			msg:   "slice marshaler error",
			value: struct{ V []*errMarshaler }{V: []*errMarshaler{{errors.New("marshal error")}}},
			err:   "marshal error",
		},
	}

	opts := test.CompareErrorMessages()

	for _, test := range tests {
		_, err := desc.Marshal(test.value)
		if diff := cmp.Diff(errors.New(test.err), err, opts); diff != "" {
			t.Errorf("%s returned unexpected error:\n%s", test.msg, diff)
		}
	}

	// Ensure write errors are passed though from each part of the application.
	// Test up to 7 writes: percent, field name, percent, newline, value, newline, newline
	for i := 1; i <= 7; i++ {
		want := errors.New("custom error")
		w := &errWriter{i, want}
		err := desc.NewEncoder(w).Encode(record{Slice: []string{"foo", "bar", "baz"}})
		if diff := cmp.Diff(err, want, cmpopts.EquateErrors()); diff != "" {
			t.Errorf("write %d should return error:\n%s", i, diff)
		}
	}
}

func BenchmarkMarshal(b *testing.B) {
	benchmarks := []struct {
		name string
		in   any
	}{
		{"string", struct{ V string }{V: "foo\nbar\nbaz"}},
		{"int", struct{ V int }{V: 1}},
		{"uint", struct{ V uint }{V: 1}},
		{"float64", struct{ V float64 }{V: 1.123}},
		{"time", struct{ V time.Time }{V: time.Date(2023, 1, 10, 19, 4, 25, 0, time.UTC)}},
		{"[8]byte", struct{ V [8]byte }{V: [8]byte{1, 2, 3, 4, 5, 6, 7, 8}}},
	}

	for _, bm := range benchmarks {
		b.Run(bm.name, func(b *testing.B) {
			for b.Loop() {
				if _, err := desc.Marshal(bm.in); err != nil {
					b.Fatal(err)
				}
			}
		})
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
