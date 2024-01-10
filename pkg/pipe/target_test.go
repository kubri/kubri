package pipe_test

import (
	"errors"
	"os"
	"path/filepath"
	"testing"

	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob/container"
	"github.com/abemedia/appcast/internal/test"
	"github.com/abemedia/appcast/pkg/pipe"
	"github.com/abemedia/appcast/target"
	"github.com/abemedia/appcast/target/azureblob"
	"github.com/abemedia/appcast/target/file"
	"github.com/abemedia/appcast/target/gcs"
	"github.com/abemedia/appcast/target/github"
	"github.com/abemedia/appcast/target/s3"
	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	gh "github.com/google/go-github/github"
	"gopkg.in/yaml.v3"
)

func TestTarget(t *testing.T) {
	dir := t.TempDir()

	tests := []struct {
		desc string
		in   string
		want func() (target.Target, error)
		err  error
	}{
		{
			desc: "file",
			in: `
				target:
					type: file
					path: ` + dir + `
			`,
			want: func() (target.Target, error) {
				return file.New(file.Config{Path: dir})
			},
		},
		{
			desc: "s3",
			in: `
				target:
					type: s3
					bucket: test
					folder: test
			`,
			want: func() (target.Target, error) {
				return s3.New(s3.Config{Bucket: "test", Folder: "test"})
			},
		},
		{
			desc: "gcs",
			in: `
				target:
					type: gcs
					bucket: test
					folder: test
			`,
			want: func() (target.Target, error) {
				t.Setenv("STORAGE_EMULATOR_HOST", "test")
				return gcs.New(gcs.Config{Bucket: "test", Folder: "test"})
			},
		},
		{
			desc: "azureblob",
			in: `
				target:
					type: azureblob
					bucket: test
					folder: test
			`,
			want: func() (target.Target, error) {
				t.Setenv("AZURE_STORAGE_ACCOUNT", "test")
				t.Setenv("AZURE_STORAGE_KEY", "test")
				return azureblob.New(azureblob.Config{Bucket: "test", Folder: "test"})
			},
		},
		{
			desc: "github",
			in: `
				target:
					type: github
					owner: abemedia
					repo: appcast
					branch: master
					folder: test
			`,
			want: func() (target.Target, error) {
				return github.New(github.Config{Owner: "abemedia", Repo: "appcast", Branch: "master", Folder: "test"})
			},
		},
		{
			desc: "invalid type",
			in: `
				target:
					type: nope
			`,
			err: errors.New("target: invalid type"),
		},
		{
			desc: "unmarshal error",
			in: `
				target:
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

		// Ignore GitHub rate limit.
		cmpopts.IgnoreTypes(gh.Rate{}),
	}

	for _, test := range tests {
		var want target.Target
		if test.want != nil {
			w, err := test.want()
			if err != nil {
				t.Fatal(err)
			}
			want = w
		}

		baseConfig := `
			source:
				type: file
				path: ` + dir + `
			apk:
				folder: .
		`

		config := append(clean(test.in), clean(baseConfig)...)
		path := filepath.Join(t.TempDir(), "appcast.yml")
		os.WriteFile(path, config, os.ModePerm)

		p, err := pipe.Load(path)

		var got target.Target
		if p != nil && p.Apk != nil {
			got = p.Apk.Target
		}

		if diff := cmp.Diff(test.err, err, opts); diff != "" {
			t.Errorf("%s:\n%s", test.desc, diff)
		} else if diff := cmp.Diff(want, got, opts); diff != "" {
			t.Errorf("%s:\n%s", test.desc, diff)
		}
	}
}
