package pipe_test

import (
	"errors"
	"os"
	"path/filepath"
	"testing"

	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob/container"
	"github.com/abemedia/appcast/internal/test"
	"github.com/abemedia/appcast/pkg/pipe"
	"github.com/abemedia/appcast/source"
	"github.com/abemedia/appcast/source/azureblob"
	"github.com/abemedia/appcast/source/file"
	"github.com/abemedia/appcast/source/gcs"
	"github.com/abemedia/appcast/source/github"
	"github.com/abemedia/appcast/source/gitlab"
	"github.com/abemedia/appcast/source/local"
	"github.com/abemedia/appcast/source/s3"
	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"gopkg.in/yaml.v3"
)

func TestSource(t *testing.T) {
	dir := t.TempDir()

	tests := []struct {
		desc string
		in   string
		want func() (*source.Source, error)
		err  error
	}{
		{
			desc: "file",
			in: `
				source:
					type: file
					path: ` + dir + `
			`,
			want: func() (*source.Source, error) {
				return file.New(file.Config{Path: dir})
			},
		},
		{
			desc: "s3",
			in: `
				source:
					type: s3
					bucket: test
					folder: test
			`,
			want: func() (*source.Source, error) {
				return s3.New(s3.Config{Bucket: "test", Folder: "test"})
			},
		},
		{
			desc: "gcs",
			in: `
				source:
					type: gcs
					bucket: test
					folder: test
			`,
			want: func() (*source.Source, error) {
				t.Setenv("STORAGE_EMULATOR_HOST", "test")
				return gcs.New(gcs.Config{Bucket: "test", Folder: "test"})
			},
		},
		{
			desc: "azureblob",
			in: `
				source:
					type: azureblob
					bucket: test
					folder: test
			`,
			want: func() (*source.Source, error) {
				t.Setenv("AZURE_STORAGE_ACCOUNT", "test")
				t.Setenv("AZURE_STORAGE_KEY", "test")
				return azureblob.New(azureblob.Config{Bucket: "test", Folder: "test"})
			},
		},
		{
			desc: "github",
			in: `
				source:
					type: github
					owner: test
					repo: test
			`,
			want: func() (*source.Source, error) {
				return github.New(github.Config{Owner: "test", Repo: "test"})
			},
		},
		{
			desc: "gitlab",
			in: `
				source:
					type: gitlab
					owner: test
					repo: test
			`,
			want: func() (*source.Source, error) {
				return gitlab.New(gitlab.Config{Owner: "test", Repo: "test"})
			},
		},
		{
			desc: "local",
			in: `
			source:
				type: local
				path: ` + dir + `
				version: v1.0.0
			`,
			want: func() (*source.Source, error) {
				return local.New(local.Config{Path: dir, Version: "v1.0.0"})
			},
		},
		{
			desc: "invalid type",
			in: `
				source:
					type: nope
			`,
			err: errors.New("source: invalid type"),
		},
		{
			desc: "unmarshal error",
			in: `
				source:
					type: {}
			`,
			err: &yaml.TypeError{Errors: []string{"line 2: cannot unmarshal !!map into string"}},
		},
	}

	opts := cmp.Options{
		test.ExportAll(),
		test.IgnoreFunctions(),
		test.CompareLoggers(),

		// Ignore azblob policies as they are not comparable.
		cmpopts.IgnoreFields(container.Client{}, "inner.internal.pl"),
	}

	for _, test := range tests {
		var want *source.Source
		if test.want != nil {
			w, err := test.want()
			if err != nil {
				t.Fatal(err)
			}
			want = w
		}

		baseConfig := `
			target:
				type: file
				path: ` + dir + `
			apk:
				folder: .
		`

		config := append(clean(test.in), clean(baseConfig)...)
		path := filepath.Join(t.TempDir(), "appcast.yml")
		os.WriteFile(path, config, os.ModePerm)

		p, err := pipe.Load(path)

		var got *source.Source
		if p != nil && p.Apk != nil {
			got = p.Apk.Source
		}

		if diff := cmp.Diff(test.err, err, opts); diff != "" {
			t.Errorf("%s:\n%s", test.desc, diff)
		} else if diff := cmp.Diff(want, got, opts); diff != "" {
			t.Errorf("%s:\n%s", test.desc, diff)
		}
	}
}
