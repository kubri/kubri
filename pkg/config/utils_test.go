package config_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/abemedia/appcast/internal/test"
	"github.com/abemedia/appcast/pkg/config"
	"github.com/google/go-cmp/cmp"
)

type testCase struct {
	desc string
	in   string
	want *config.Config
	err  error
	hook func()
}

func runTest(t *testing.T, tests []testCase) {
	t.Helper()

	opts := cmp.Options{
		test.ExportAll(),
		test.ComparePGPKeys(),
		test.CompareRSAPrivateKeys(),
	}

	for _, tc := range tests {
		t.Setenv("APPCAST_PATH", t.TempDir())

		if tc.hook != nil {
			tc.hook()
		}

		path := filepath.Join(t.TempDir(), "appcast.yml")
		os.WriteFile(path, test.YAML(tc.in), os.ModePerm)

		got, err := config.Load(path)

		if diff := cmp.Diff(tc.err, err, opts); diff != "" {
			t.Errorf("%s:\n%s", tc.desc, diff)
		} else if diff := cmp.Diff(tc.want, got, opts); diff != "" {
			t.Errorf("%s:\n%s", tc.desc, diff)
		}
	}
}
