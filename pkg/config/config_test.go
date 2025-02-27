package config_test

import (
	"errors"
	"io/fs"
	"os"
	"path/filepath"
	"testing"

	"github.com/google/go-cmp/cmp"
	"gopkg.in/yaml.v3"

	"github.com/kubri/kubri/internal/test"
	"github.com/kubri/kubri/pkg/config"
)

func TestLoad(t *testing.T) {
	tests := []struct {
		desc   string
		config string
		path   string
		mode   fs.FileMode
		err    error
	}{
		{path: "kubri.yml"},
		{path: "kubri.yaml"},
		{path: ".kubri.yml"},
		{path: ".kubri.yaml"},
		{path: filepath.Join(".github", "kubri.yml")},
		{path: filepath.Join(".github", "kubri.yaml")},
		{
			path: "foo.yml",
			err:  errors.New("no config file found"),
		},
		{
			desc: "permission denied",
			path: "kubri.yml",
			err:  errors.New("open kubri.yml: permission denied"),
			mode: 0o200,
		},
		{
			desc:   "invalid yaml",
			config: `*&%^`,
			path:   "kubri.yml",
			err:    errors.New("yaml: did not find expected alphabetic or numeric character"),
		},
		{
			desc:   "non-existent field",
			config: `foo: bar`,
			path:   "kubri.yml",
			err:    &yaml.TypeError{Errors: []string{"line 1: field foo not found in type config.config"}},
		},
		{
			desc:   "failed validation",
			config: `version: invalid`,
			path:   "kubri.yml",
			err: &config.Error{
				Errors: []string{
					"version must be a valid version constraint",
					"source is a required field",
					"target is a required field",
				},
			},
		},
	}

	opts := cmp.Options{
		test.ExportAll(),
		test.CompareErrorMessages(),
	}

	for _, tc := range tests {
		if tc.desc == "" {
			tc.desc = tc.path
		}
		if tc.config == "" {
			tc.config = `
				source:
					type: file
					path: ` + t.TempDir() + `
				target:
					type: file
					path: ` + t.TempDir()
		}
		if tc.mode == 0 {
			tc.mode = os.ModePerm
		}

		t.Chdir(t.TempDir())
		os.MkdirAll(filepath.Dir(tc.path), os.ModePerm)
		os.WriteFile(tc.path, test.YAML(tc.config), tc.mode)

		_, err := config.Load("")

		if diff := cmp.Diff(tc.err, err, opts); diff != "" {
			t.Errorf("%s:\n%s", tc.path, diff)
		}
	}
}
