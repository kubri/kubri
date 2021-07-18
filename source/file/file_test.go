package file_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/abemedia/appcast/source"
	"github.com/abemedia/appcast/source/file"
	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
)

func TestFile(t *testing.T) {
	dir := t.TempDir()

	want := []*source.Release{
		{
			Name:    "v1.0.0",
			Version: "v1.0.0",
			Assets: []*source.Asset{
				{
					Name: "test.dmg",
					URL:  "file://" + filepath.Join(dir, "v1.0.0", "test.dmg"),
					Size: 5,
				},
				{
					Name: "test_32-bit.msi",
					URL:  "file://" + filepath.Join(dir, "v1.0.0", "test_32-bit.msi"),
					Size: 5,
				},
				{
					Name: "test_64-bit.msi",
					URL:  "file://" + filepath.Join(dir, "v1.0.0", "test_64-bit.msi"),
					Size: 5,
				},
			},
		},
		{
			Name:       "v1.0.0-alpha1",
			Version:    "v1.0.0-alpha1",
			Prerelease: true,
			Assets: []*source.Asset{
				{
					Name: "test.dmg",
					URL:  "file://" + filepath.Join(dir, "v1.0.0-alpha1", "test.dmg"),
					Size: 5,
				},
				{
					Name: "test_32-bit.msi",
					URL:  "file://" + filepath.Join(dir, "v1.0.0-alpha1", "test_32-bit.msi"),
					Size: 5,
				},
				{
					Name: "test_64-bit.msi",
					URL:  "file://" + filepath.Join(dir, "v1.0.0-alpha1", "test_64-bit.msi"),
					Size: 5,
				},
			},
		},
	}

	for _, release := range want {
		path := filepath.Join(dir, release.Version)
		err := os.Mkdir(path, os.ModePerm)
		if err != nil {
			t.Fatal(err)
		}

		for _, asset := range release.Assets {
			err = os.WriteFile(filepath.Join(path, asset.Name), []byte("test\n"), os.ModePerm)
			if err != nil {
				t.Fatal(err)
			}
		}
	}

	r, err := file.New(source.Config{Repo: dir})
	if err != nil {
		t.Fatal(err)
	}

	got, err := r.Releases()
	if err != nil {
		t.Fatal(err)
	}

	if diff := cmp.Diff(want, got, cmpopts.IgnoreFields(source.Release{}, "Date")); diff != "" {
		t.Error(diff)
	}
}
