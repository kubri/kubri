package pipe //nolint:testpackage

import (
	"testing"

	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob/container"
	"github.com/abemedia/appcast/internal/test"
	"github.com/abemedia/appcast/target"
	"github.com/abemedia/appcast/target/azureblob"
	"github.com/abemedia/appcast/target/file"
	"github.com/abemedia/appcast/target/gcs"
	"github.com/abemedia/appcast/target/github"
	"github.com/abemedia/appcast/target/s3"
	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	gh "github.com/google/go-github/github"
	"github.com/mitchellh/mapstructure"
)

func TestTarget(t *testing.T) {
	dir := t.TempDir()

	tests := []struct {
		in   targetConfig
		want func() (target.Target, error)
		err  error
	}{
		{
			in: targetConfig{
				"type": "file",
				"path": dir,
			},
			want: func() (target.Target, error) {
				return file.New(file.Config{Path: dir})
			},
		},
		{
			in: targetConfig{
				"type": "file",
				"path": 1,
			},
			err: &mapstructure.Error{Errors: []string{"'Path' expected type 'string', got unconvertible type 'int', value: '1'"}},
		},
		{
			in: targetConfig{
				"type":   "s3",
				"bucket": "test",
				"folder": "test",
			},
			want: func() (target.Target, error) {
				return s3.New(s3.Config{Bucket: "test", Folder: "test"})
			},
		},
		{
			in: targetConfig{
				"type":   "s3",
				"bucket": 1,
			},
			err: &mapstructure.Error{Errors: []string{"'Bucket' expected type 'string', got unconvertible type 'int', value: '1'"}},
		},
		{
			in: targetConfig{
				"type":   "gcs",
				"bucket": "test",
				"folder": "test",
			},
			want: func() (target.Target, error) {
				t.Setenv("STORAGE_EMULATOR_HOST", "test")
				return gcs.New(gcs.Config{Bucket: "test", Folder: "test"})
			},
		},
		{
			in: targetConfig{
				"type":   "gcs",
				"bucket": 1,
			},
			err: &mapstructure.Error{Errors: []string{"'Bucket' expected type 'string', got unconvertible type 'int', value: '1'"}},
		},
		{
			in: targetConfig{
				"type":   "azureblob",
				"bucket": "test",
				"folder": "test",
			},
			want: func() (target.Target, error) {
				t.Setenv("AZURE_STORAGE_ACCOUNT", "test")
				t.Setenv("AZURE_STORAGE_KEY", "test")
				return azureblob.New(azureblob.Config{Bucket: "test", Folder: "test"})
			},
		},
		{
			in: targetConfig{
				"type":   "azureblob",
				"bucket": 1,
			},
			err: &mapstructure.Error{Errors: []string{"'Bucket' expected type 'string', got unconvertible type 'int', value: '1'"}},
		},
		{
			in: targetConfig{
				"type":  "github",
				"owner": "abemedia",
				"repo":  "appcast",
			},
			want: func() (target.Target, error) {
				return github.New(github.Config{Owner: "abemedia", Repo: "appcast"})
			},
		},
		{
			in: targetConfig{
				"type":  "github",
				"owner": 1,
			},
			err: &mapstructure.Error{Errors: []string{"'Owner' expected type 'string', got unconvertible type 'int', value: '1'"}},
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

		got, err := getTarget(test.in)

		if diff := cmp.Diff(test.err, err, opts); diff != "" {
			t.Errorf("%s:\n%s", test.in["type"], diff)
		} else if diff := cmp.Diff(want, got, opts); diff != "" {
			t.Errorf("%s:\n%s", test.in["type"], diff)
		}
	}
}
