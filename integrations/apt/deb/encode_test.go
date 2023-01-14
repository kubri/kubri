package deb_test

import (
	"testing"
	"time"

	"github.com/abemedia/appcast/integrations/apt/deb"
	"github.com/google/go-cmp/cmp"
)

func TestMarshal(t *testing.T) {
	in := record{
		String:    "test",
		ByteArray: [4]byte{1, 2, 3, 4},
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
		in   interface{}
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
			msg: "multi-line string",
			in: record{
				String: "foo\nbar\nbaz\n\nfoobar",
			},
			want: `String: foo
 bar
 baz
 .
 foobar
`,
		},
		{
			msg: "multi-line string starting with empty line",
			in: record{
				String: "\nfoo\nbar",
			},
			want: `String:
 foo
 bar
`,
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

func BenchmarkMarshal(b *testing.B) {
	type record struct {
		String    string
		ByteArray [4]byte `deb:"Hex"`
		Int       int
	}

	v := record{
		String:    "test",
		ByteArray: [4]byte{1, 2, 3, 4},
		Int:       1,
	}

	for i := 0; i < b.N; i++ {
		deb.Marshal(v)
	}
}
