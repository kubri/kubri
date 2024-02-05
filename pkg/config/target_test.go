package config_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob/container"
	"github.com/abemedia/appcast/internal/test"
	"github.com/abemedia/appcast/pkg/config"
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
		desc   string
		config string
		want   func() (target.Target, error)
		err    error
	}{
		{
			desc: "file",
			config: `
				target:
					type: file
					path: ` + dir + `
			`,
			want: func() (target.Target, error) {
				return file.New(file.Config{Path: dir})
			},
		},
		{
			desc: "file invalid",
			config: `
				target:
					type: file
					url: invalid
			`,
			err: &config.Error{
				Errors: []string{
					"target.path is a required field",
					"target.url must be a valid URL",
				},
			},
		},
		{
			desc: "s3",
			config: `
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
			desc: "s3 invalid",
			config: `
				target:
					type: s3
					folder: '*'
					endpoint: invalid
					url: invalid
			`,
			err: &config.Error{
				Errors: []string{
					"target.bucket is a required field",
					"target.folder must be a valid folder name",
					"target.endpoint must be a valid URL or FQDN",
					"target.url must be a valid URL",
				},
			},
		},
		{
			desc: "gcs",
			config: `
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
			desc: "gcs invalid",
			config: `
				target:
					type: gcs
					folder: '*'
					url: invalid
			`,
			err: &config.Error{
				Errors: []string{
					"target.bucket is a required field",
					"target.folder must be a valid folder name",
					"target.url must be a valid URL",
				},
			},
		},
		{
			desc: "azureblob",
			config: `
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
			desc: "azureblob invalid",
			config: `
				target:
					type: azureblob
					folder: '*'
					url: invalid
			`,
			err: &config.Error{
				Errors: []string{
					"target.bucket is a required field",
					"target.folder must be a valid folder name",
					"target.url must be a valid URL",
				},
			},
		},
		{
			desc: "github",
			config: `
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
			desc: "github invalid",
			config: `
				target:
					type: github
					folder: '*'
			`,
			err: &config.Error{
				Errors: []string{
					"target.owner is a required field",
					"target.repo is a required field",
					"target.folder must be a valid folder name",
				},
			},
		},
		{
			desc: "invalid type",
			config: `
				target:
					type: nope
			`,
			err: &config.Error{Errors: []string{"target.type must be one of [azureblob gcs s3 file github]"}},
		},
		{
			desc: "unmarshal error",
			config: `
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

	for _, tc := range tests {
		var want target.Target
		if tc.want != nil {
			w, err := tc.want()
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

		path := filepath.Join(t.TempDir(), "appcast.yml")
		os.WriteFile(path, test.JoinYAML(tc.config, baseConfig), os.ModePerm)

		p, err := config.Load(path)

		var got target.Target
		if p != nil && p.Apk != nil {
			got = p.Apk.Target
		}

		if diff := cmp.Diff(tc.err, err, opts); diff != "" {
			t.Errorf("%s:\n%s", tc.desc, diff)
		} else if diff := cmp.Diff(want, got, opts); diff != "" {
			t.Errorf("%s:\n%s", tc.desc, diff)
		}
	}
}
