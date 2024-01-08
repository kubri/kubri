package pipe //nolint:testpackage

import (
	"testing"

	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob/container"
	"github.com/abemedia/appcast/internal/test"
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
	"github.com/mitchellh/mapstructure"
)

func TestSource(t *testing.T) {
	dir := t.TempDir()

	tests := []struct {
		in   sourceConfig
		want func() (*source.Source, error)
		err  error
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
				"type": "file",
				"path": 1,
			},
			err: &mapstructure.Error{Errors: []string{"'Path' expected type 'string', got unconvertible type 'int', value: '1'"}},
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
				"type":   "s3",
				"bucket": 1,
			},
			err: &mapstructure.Error{Errors: []string{"'Bucket' expected type 'string', got unconvertible type 'int', value: '1'"}},
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
				"type":   "gcs",
				"bucket": 1,
			},
			err: &mapstructure.Error{Errors: []string{"'Bucket' expected type 'string', got unconvertible type 'int', value: '1'"}},
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
				"type":   "azureblob",
				"bucket": 1,
			},
			err: &mapstructure.Error{Errors: []string{"'Bucket' expected type 'string', got unconvertible type 'int', value: '1'"}},
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
				"type":  "github",
				"owner": 1,
			},
			err: &mapstructure.Error{Errors: []string{"'Owner' expected type 'string', got unconvertible type 'int', value: '1'"}},
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
				"type":  "gitlab",
				"owner": 1,
			},
			err: &mapstructure.Error{Errors: []string{"'Owner' expected type 'string', got unconvertible type 'int', value: '1'"}},
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
		{
			in: sourceConfig{
				"type": "local",
				"path": 1,
			},
			err: &mapstructure.Error{Errors: []string{"'Path' expected type 'string', got unconvertible type 'int', value: '1'"}},
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

		got, err := getSource(test.in)

		if diff := cmp.Diff(test.err, err, opts); diff != "" {
			t.Errorf("%s:\n%s", test.in["type"], diff)
		} else if diff := cmp.Diff(want, got, opts); diff != "" {
			t.Errorf("%s:\n%s", test.in["type"], diff)
		}
	}
}
