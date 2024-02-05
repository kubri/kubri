package config_test

import (
	"errors"
	"io/fs"
	"os"
	"path/filepath"
	"testing"

	"github.com/abemedia/appcast/internal/test"
	"github.com/abemedia/appcast/pkg/config"
	"github.com/google/go-cmp/cmp"
	"gopkg.in/yaml.v3"
)

func TestLoad(t *testing.T) {
	tests := []struct {
		desc   string
		config string
		path   string
		mode   fs.FileMode
		err    error
	}{
		{path: "appcast.yml"},
		{path: "appcast.yaml"},
		{path: ".appcast.yml"},
		{path: ".appcast.yaml"},
		{path: filepath.Join(".github", "appcast.yml")},
		{path: filepath.Join(".github", "appcast.yaml")},
		{
			path: "foo.yml",
			err:  errors.New("no config file found"),
		},
		{
			desc: "permission denied",
			path: "appcast.yml",
			err:  errors.New("open appcast.yml: permission denied"),
			mode: 0o200,
		},
		{
			desc:   "invalid yaml",
			config: `*&%^`,
			path:   "appcast.yml",
			err:    errors.New("yaml: did not find expected alphabetic or numeric character"),
		},
		{
			desc:   "non-existent field",
			config: `foo: bar`,
			path:   "appcast.yml",
			err:    &yaml.TypeError{Errors: []string{"line 1: field foo not found in type config.config"}},
		},
		{
			desc:   "failed validation",
			config: `version: invalid`,
			path:   "appcast.yml",
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

	wd, _ := os.Getwd()
	defer os.Chdir(wd)

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

		os.Chdir(t.TempDir())
		os.MkdirAll(filepath.Dir(tc.path), os.ModePerm)
		os.WriteFile(tc.path, test.YAML(tc.config), tc.mode)

		_, err := config.Load("")

		if diff := cmp.Diff(tc.err, err, opts); diff != "" {
			t.Errorf("%s:\n%s", tc.path, diff)
		}
	}
}
