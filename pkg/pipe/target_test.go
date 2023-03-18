package pipe //nolint:testpackage

import (
	"reflect"
	"testing"

	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob/container"
	"github.com/abemedia/appcast/target"
	"github.com/abemedia/appcast/target/blob/azureblob"
	"github.com/abemedia/appcast/target/blob/gcs"
	"github.com/abemedia/appcast/target/blob/s3"
	"github.com/abemedia/appcast/target/file"
	"github.com/abemedia/appcast/target/github"
	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	gh "github.com/google/go-github/github"
)

func TestTarget(t *testing.T) {
	dir := t.TempDir()

	tests := []struct {
		in   targetConfig
		want func() (target.Target, error)
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
				"type":   "gcs",
				"bucket": "test",
				"folder": "test",
			},
			want: func() (target.Target, error) {
				return gcs.New(gcs.Config{Bucket: "test", Folder: "test"})
			},
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
				"type":  "github",
				"owner": "abemedia",
				"repo":  "appcast-test",
			},
			want: func() (target.Target, error) {
				return github.New(github.Config{Owner: "abemedia", Repo: "appcast-test"})
			},
		},
	}

	for i, test := range tests {
		want, err := test.want()
		if err != nil {
			t.Fatal(err)
		}

		s, err := getTarget(test.in)
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

			// Ignore GitHub rate limit.
			cmpopts.IgnoreTypes(gh.Rate{}),
		}

		if diff := cmp.Diff(want, s, opts...); diff != "" {
			t.Error(i, diff)
		}
	}
}
