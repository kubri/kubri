package config_test

import (
	"errors"
	"io/fs"
	"os"
	"path/filepath"
	"syscall"
	"testing"

	"github.com/goccy/go-yaml"
	"github.com/goccy/go-yaml/token"
	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"

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
			err:  &fs.PathError{Op: "open", Path: "kubri.yml", Err: syscall.EACCES},
			mode: 0o200,
		},
		{
			desc:   "invalid yaml",
			config: `*&%^`,
			path:   "kubri.yml",
			err:    &yaml.SyntaxError{Message: "unexpected directive value. document not started"},
		},
		{
			desc:   "non-existent field",
			config: `foo: bar`,
			path:   "kubri.yml",
			err:    &yaml.UnknownFieldError{Message: `unknown field "foo"`},
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
		cmpopts.IgnoreTypes(&token.Token{}),
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
