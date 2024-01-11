package pipe_test

import (
	"context"
	"errors"
	"io/fs"
	"os"
	"path/filepath"
	"testing"

	"github.com/abemedia/appcast/integrations/apk"
	"github.com/abemedia/appcast/integrations/appinstaller"
	"github.com/abemedia/appcast/integrations/apt"
	"github.com/abemedia/appcast/integrations/sparkle"
	"github.com/abemedia/appcast/integrations/yum"
	"github.com/abemedia/appcast/internal/test"
	"github.com/abemedia/appcast/internal/testsource"
	"github.com/abemedia/appcast/pkg/pipe"
	"github.com/abemedia/appcast/source"
	target "github.com/abemedia/appcast/target/file"
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
			err:    &yaml.TypeError{Errors: []string{"line 1: field foo not found in type pipe.config"}},
		},
		{
			desc:   "failed validation",
			config: `version: invalid`,
			path:   "appcast.yml",
			err: &pipe.Error{
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

	for _, test := range tests {
		if test.desc == "" {
			test.desc = test.path
		}
		if test.config == "" {
			test.config = `
				source:
					type: file
					path: ` + t.TempDir() + `
				target:
					type: file
					path: ` + t.TempDir()
		}
		if test.mode == 0 {
			test.mode = os.ModePerm
		}

		os.Chdir(t.TempDir())
		os.MkdirAll(filepath.Dir(test.path), os.ModePerm)
		os.WriteFile(test.path, clean(test.config), test.mode)

		_, err := pipe.Load("")

		if diff := cmp.Diff(test.err, err, opts); diff != "" {
			t.Errorf("%s:\n%s", test.path, diff)
		}
	}
}

func TestPipe(t *testing.T) {
	dir := t.TempDir()
	src := testsource.New([]*source.Release{{Version: "v1.0.0"}})
	tgt, _ := target.New(target.Config{Path: dir})

	tests := []struct {
		desc string
		pipe *pipe.Pipe
		err  error
	}{
		{
			desc: "empty",
			pipe: &pipe.Pipe{},
			err:  errors.New("no integrations configured"),
		},
		{
			desc: "all",
			pipe: &pipe.Pipe{
				Apk: &apk.Config{
					Source: src,
					Target: tgt,
				},
				Appinstaller: &appinstaller.Config{
					Source: src,
					Target: tgt,
				},
				Apt: &apt.Config{
					Source: src,
					Target: tgt,
				},
				Sparkle: &sparkle.Config{
					FileName: "appcast.xml",
					Source:   src,
					Target:   tgt,
				},
				Yum: &yum.Config{
					Source: src,
					Target: tgt,
				},
			},
		},
		{
			desc: "appinstaller error",
			pipe: &pipe.Pipe{
				Appinstaller: &appinstaller.Config{
					Source: source.New(nil),
					Target: tgt,
				},
			},
			err: errors.New("failed to publish App Installer packages: missing source"),
		},
		{
			desc: "apk error",
			pipe: &pipe.Pipe{
				Apk: &apk.Config{
					Source: source.New(nil),
					Target: tgt,
				},
			},
			err: errors.New("failed to publish APK packages: missing source"),
		},
		{
			desc: "apt error",
			pipe: &pipe.Pipe{
				Apt: &apt.Config{
					Source: source.New(nil),
					Target: tgt,
				},
			},
			err: errors.New("failed to publish APT packages: missing source"),
		},
		{
			desc: "sparkle error",
			pipe: &pipe.Pipe{
				Sparkle: &sparkle.Config{
					FileName: "appcast.xml",
					Source:   source.New(nil),
					Target:   tgt,
				},
			},
			err: errors.New("failed to publish Sparkle packages: missing source"),
		},
		{
			desc: "yum error",
			pipe: &pipe.Pipe{
				Yum: &yum.Config{
					Source: source.New(nil),
					Target: tgt,
				},
			},
			err: errors.New("failed to publish YUM packages: missing source"),
		},
	}

	opts := cmp.Options{
		test.ExportAll(),
		test.CompareErrorMessages(),
	}

	for _, test := range tests {
		err := test.pipe.Run(context.Background())
		if diff := cmp.Diff(test.err, err, opts); diff != "" {
			t.Errorf("%s:\n%s", test.desc, diff)
		}
	}
}
