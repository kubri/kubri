package memory_test

import (
	"bytes"
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

	opts := []cmp.Option{
		cmpopts.IgnoreFields(source.Release{}, "Date"),
		cmpopts.SortSlices(func(a, b *source.Release) bool { return a.Version > b.Version }),
		cmpopts.SortSlices(func(a, b *source.Asset) bool { return a.Name > b.Name }),
	}

	t.Run("ListReleases", func(t *testing.T) {
		got, err := r.ListReleases(nil)
		if err != nil {
			t.Fatal(err)
		}

		if diff := cmp.Diff(want, got, opts...); diff != "" {
			t.Error(diff)
		}
	})

	t.Run("GetRelease", func(t *testing.T) {
		got, err := r.GetRelease(want[0].Version)
		if err != nil {
			t.Fatal(err)
		}

		if diff := cmp.Diff(want[0], got, opts...); diff != "" {
			t.Error(diff)
		}
	})

	asset := []byte("foo")

	t.Run("UploadAsset", func(t *testing.T) {
		err := r.UploadAsset(want[0].Version, "test.txt", asset)
		if err != nil {
			t.Fatal(err)
		}
	})

	t.Run("DownloadAsset", func(t *testing.T) {
		b, err := r.DownloadAsset(want[0].Version, "test.txt")
		if err != nil {
			t.Fatal(err)
		}

		if !bytes.Equal(asset, b) {
			t.Error("should be equal")
		}
	})
}
