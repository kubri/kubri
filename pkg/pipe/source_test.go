package pipe //nolint:testpackage

import (
	"reflect"
	"testing"

	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob/container"
	"github.com/abemedia/appcast/source"
	"github.com/abemedia/appcast/source/blob/azureblob"
	"github.com/abemedia/appcast/source/blob/file"
	"github.com/abemedia/appcast/source/blob/gcs"
	"github.com/abemedia/appcast/source/blob/s3"
	"github.com/abemedia/appcast/source/github"
	"github.com/abemedia/appcast/source/gitlab"
	"github.com/abemedia/appcast/source/local"
	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
)

func TestSource(t *testing.T) {
	dir := t.TempDir()

	tests := []struct {
		in   sourceConfig
		want func() (*source.Source, error)
	}{
		{
			in: sourceConfig{
				"type": "file",
				"path": dir,
			},
			want: func() (*source.Source, error) {
				return file.New(file.Config{Path: dir})
			},
		},
		{
			in: sourceConfig{
				"type":   "s3",
				"bucket": "test",
				"folder": "test",
			},
			want: func() (*source.Source, error) {
				return s3.New(s3.Config{Bucket: "test", Folder: "test"})
			},
		},
		{
			in: sourceConfig{
				"type":   "gcs",
				"bucket": "test",
				"folder": "test",
			},
			want: func() (*source.Source, error) {
				t.Setenv("STORAGE_EMULATOR_HOST", "test")
				return gcs.New(gcs.Config{Bucket: "test", Folder: "test"})
			},
		},
		{
			in: sourceConfig{
				"type":   "azureblob",
				"bucket": "test",
				"folder": "test",
			},
			want: func() (*source.Source, error) {
				t.Setenv("AZURE_STORAGE_ACCOUNT", "test")
				t.Setenv("AZURE_STORAGE_KEY", "test")
				return azureblob.New(azureblob.Config{Bucket: "test", Folder: "test"})
			},
		},
		{
			in: sourceConfig{
				"type":  "github",
				"owner": "test",
				"repo":  "test",
			},
			want: func() (*source.Source, error) {
				return github.New(github.Config{Owner: "test", Repo: "test"})
			},
		},
		{
			in: sourceConfig{
				"type":  "gitlab",
				"owner": "test",
				"repo":  "test",
			},
			want: func() (*source.Source, error) {
				return gitlab.New(gitlab.Config{Owner: "test", Repo: "test"})
			},
		},
		{
			in: sourceConfig{
				"type":    "local",
				"path":    dir,
				"version": "v1.0.0",
			},
			want: func() (*source.Source, error) {
				return local.New(local.Config{Path: dir, Version: "v1.0.0"})
			},
		},
	}

	for i, test := range tests {
		want, err := test.want()
		if err != nil {
			t.Fatal(err)
		}

		s, err := getSource(test.in)
		if err != nil {
			t.Error(err)
			continue
		}

		opts := cmp.Options{
			// Export all unexported fields.
			cmp.Exporter(func(t reflect.Type) bool { return true }),

			// Ignore all function fields.
			cmp.FilterPath(func(p cmp.Path) bool {
				sf, ok := p.Index(-1).(cmp.StructField)
				return ok && sf.Type().Kind() == reflect.Func
			}, cmp.Ignore()),

			// Ignore azblob policies as they are not comparable.
			cmpopts.IgnoreFields(container.Client{}, "inner.pl"),
		}

		if diff := cmp.Diff(want, s, opts...); diff != "" {
			t.Error(i, diff)
		}
	}
}
