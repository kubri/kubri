package config_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/google/go-cmp/cmp"

	"github.com/kubri/kubri/internal/test"
	"github.com/kubri/kubri/pkg/config"
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
		test.IgnoreFunctions(),
	}

	for _, tc := range tests {
		t.Setenv("KUBRI_PATH", t.TempDir())

		if tc.hook != nil {
			tc.hook()
		}

		path := filepath.Join(t.TempDir(), "kubri.yml")
		os.WriteFile(path, test.YAML(tc.in), os.ModePerm)

		got, err := config.Load(path)

		if diff := cmp.Diff(tc.err, err, opts); diff != "" {
			t.Errorf("%s:\n%s", tc.desc, diff)
		} else if diff := cmp.Diff(tc.want, got, opts); diff != "" {
			t.Errorf("%s:\n%s", tc.desc, diff)
		}
	}
}
