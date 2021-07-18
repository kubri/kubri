package memory_test

import (
	"testing"

	"github.com/abemedia/appcast/source"
	"github.com/abemedia/appcast/source/memory"
	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
)

func TestLocal(t *testing.T) {
	want := []*source.Release{
		{
			Name:    "v1.0.0",
			Version: "v1.0.0",
			Assets: []*source.Asset{
				{
					Name: "test.dmg",
					URL:  "memory://v1.0.0/test.dmg",
					Size: 5,
				},
				{
					Name: "test_32-bit.msi",
					URL:  "memory://v1.0.0/test_32-bit.msi",
					Size: 5,
				},
				{
					Name: "test_64-bit.msi",
					URL:  "memory://v1.0.0/test_64-bit.msi",
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
					URL:  "memory://v1.0.0-alpha1/test.dmg",
					Size: 5,
				},
				{
					Name: "test_32-bit.msi",
					URL:  "memory://v1.0.0-alpha1/test_32-bit.msi",
					Size: 5,
				},
				{
					Name: "test_64-bit.msi",
					URL:  "memory://v1.0.0-alpha1/test_64-bit.msi",
					Size: 5,
				},
			},
		},
	}

	r, _ := memory.New(source.Config{})

	for _, release := range want {
		for _, asset := range release.Assets {
			r.UploadAsset(release.Version, asset.Name, []byte("test\n"))
		}
	}

	got, err := r.Releases()
	if err != nil {
		t.Fatal(err)
	}

	if diff := cmp.Diff(want, got, cmpopts.IgnoreFields(source.Release{}, "Date")); diff != "" {
		t.Error(diff)
	}
}
