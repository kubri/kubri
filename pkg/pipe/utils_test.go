package pipe_test

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/MakeNowJust/heredoc/v2"
	"github.com/abemedia/appcast/internal/test"
	"github.com/abemedia/appcast/pkg/pipe"
	"github.com/google/go-cmp/cmp"
)

type testCase struct {
	desc string
	in   string
	want *pipe.Pipe
	err  error
	hook func()
}

func clean(s string) []byte {
	return []byte(strings.ReplaceAll(heredoc.Doc(s), "\t", "  "))
}

func runTest(t *testing.T, tests []testCase) {
	t.Helper()

	opts := cmp.Options{
		test.ExportAll(),
		test.ComparePGPKeys(),
		test.CompareRSAPrivateKeys(),
	}

	for _, test := range tests {
		t.Setenv("APPCAST_PATH", t.TempDir())

		if test.hook != nil {
			test.hook()
		}

		path := filepath.Join(t.TempDir(), "appcast.yml")
		os.WriteFile(path, clean(test.in), os.ModePerm)

		got, err := pipe.Load(path)

		if diff := cmp.Diff(test.err, err, opts); diff != "" {
			t.Errorf("%s:\n%s", test.desc, diff)
		} else if diff := cmp.Diff(test.want, got, opts); diff != "" {
			t.Errorf("%s:\n%s", test.desc, diff)
		}
	}
}
