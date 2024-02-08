package config_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob/container"
	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"gopkg.in/yaml.v3"

	"github.com/kubri/kubri/internal/test"
	"github.com/kubri/kubri/pkg/config"
	"github.com/kubri/kubri/source"
	"github.com/kubri/kubri/source/azureblob"
	"github.com/kubri/kubri/source/file"
	"github.com/kubri/kubri/source/gcs"
	"github.com/kubri/kubri/source/github"
	"github.com/kubri/kubri/source/gitlab"
	"github.com/kubri/kubri/source/local"
	"github.com/kubri/kubri/source/s3"
)

func TestSource(t *testing.T) {
	dir := t.TempDir()

	tests := []struct {
		desc   string
		config string
		want   func() (*source.Source, error)
		err    error
	}{
		{
			desc: "file",
			config: `
				source:
					type: file
					path: ` + dir + `
			`,
			want: func() (*source.Source, error) {
				return file.New(file.Config{Path: dir})
			},
		},
		{
			desc: "file invalid",
			config: `
				source:
					type: file
					path: nope
					url: invalid
			`,
			err: &config.Error{
				Errors: []string{
					"source.path must be a valid path to a directory",
					"source.url must be a valid URL",
				},
			},
		},
		{
			desc: "file missing path",
			config: `
				source:
					type: file
			`,
			err: &config.Error{Errors: []string{"source.path is a required field"}},
		},
		{
			desc: "s3",
			config: `
				source:
					type: s3
					bucket: test
					folder: test
					endpoint: s3.example.com
					region: auto
					disable-ssl: true
					url: http://example.com
			`,
			want: func() (*source.Source, error) {
				return s3.New(s3.Config{
					Bucket:     "test",
					Folder:     "test",
					Endpoint:   "s3.example.com",
					Region:     "auto",
					DisableSSL: true,
					URL:        "http://example.com",
				})
			},
		},
		{
			desc: "s3 invalid",
			config: `
				source:
					type: s3
					folder: '*'
					endpoint: invalid
					url: invalid
			`,
			err: &config.Error{
				Errors: []string{
					"source.bucket is a required field",
					"source.folder must be a valid folder name",
					"source.endpoint must be a valid URL or FQDN",
					"source.url must be a valid URL",
				},
			},
		},
		{
			desc: "gcs",
			config: `
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
			desc: "gcs invalid",
			config: `
				source:
					type: gcs
					folder: '*'
					url: invalid
			`,
			err: &config.Error{
				Errors: []string{
					"source.bucket is a required field",
					"source.folder must be a valid folder name",
					"source.url must be a valid URL",
				},
			},
		},
		{
			desc: "azureblob",
			config: `
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
			desc: "azureblob invalid",
			config: `
				source:
					type: azureblob
					folder: '*'
					url: invalid
			`,
			err: &config.Error{
				Errors: []string{
					"source.bucket is a required field",
					"source.folder must be a valid folder name",
					"source.url must be a valid URL",
				},
			},
		},
		{
			desc: "github",
			config: `
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
			desc: "github invalid",
			config: `
				source:
					type: github
			`,
			err: &config.Error{
				Errors: []string{
					"source.owner is a required field",
					"source.repo is a required field",
				},
			},
		},
		{
			desc: "gitlab",
			config: `
				source:
					type: gitlab
					owner: test
					repo: test
					url: http://example.com
			`,
			want: func() (*source.Source, error) {
				return gitlab.New(gitlab.Config{Owner: "test", Repo: "test", URL: "http://example.com"})
			},
		},
		{
			desc: "gitlab invalid",
			config: `
				source:
					type: gitlab
					url: invalid
			`,
			err: &config.Error{
				Errors: []string{
					"source.owner is a required field",
					"source.repo is a required field",
					"source.url must be a valid URL",
				},
			},
		},
		{
			desc: "local",
			config: `
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
			desc: "local invalid",
			config: `
				source:
					type: local
					path: nope
					version: invalid
			`,
			err: &config.Error{
				Errors: []string{
					"source.path must be a valid path to a directory",
					"source.version must be a valid semver version",
				},
			},
		},
		{
			desc: "local missing fields",
			config: `
				source:
					type: local
			`,
			err: &config.Error{
				Errors: []string{
					"source.path is a required field",
					"source.version is a required field",
				},
			},
		},
		{
			desc: "invalid type",
			config: `
				source:
					type: nope
			`,
			err: &config.Error{Errors: []string{"source.type must be one of [azureblob gcs s3 file github gitlab local]"}},
		},
		{
			desc: "unmarshal error",
			config: `
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

	for _, tc := range tests {
		var want *source.Source
		if tc.want != nil {
			w, err := tc.want()
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

		path := filepath.Join(t.TempDir(), "kubri.yml")
		os.WriteFile(path, test.JoinYAML(tc.config, baseConfig), os.ModePerm)

		p, err := config.Load(path)

		var got *source.Source
		if p != nil && p.Apk != nil {
			got = p.Apk.Source
		}

		if diff := cmp.Diff(tc.err, err, opts); diff != "" {
			t.Errorf("%s:\n%s", tc.desc, diff)
		} else if diff := cmp.Diff(want, got, opts); diff != "" {
			t.Errorf("%s:\n%s", tc.desc, diff)
		}
	}
}
