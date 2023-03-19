package test

import (
	"bytes"
	"context"
	"testing"
	"time"

	"github.com/abemedia/appcast/source"
	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
)

//nolint:funlen
func Source(t *testing.T, s *source.Source, makeURL func(version, asset string) string) {
	t.Helper()

	data := []byte("test\n")
	ctx := context.Background()

	want := []*source.Release{
		{
			Name:    "v1.0.0",
			Date:    time.Now(),
			Version: "v1.0.0",
			Assets: []*source.Asset{
				{Name: "test.dmg", Size: 5},
				{Name: "test_32-bit.msi", Size: 5},
				{Name: "test_64-bit.msi", Size: 5},
			},
		},
		{
			Name:       "v1.0.0-alpha1",
			Date:       time.Now(),
			Version:    "v1.0.0-alpha1",
			Prerelease: true,
			Assets: []*source.Asset{
				{Name: "test.dmg", Size: 5},
				{Name: "test_32-bit.msi", Size: 5},
				{Name: "test_64-bit.msi", Size: 5},
			},
		},
		{
			Name:       "v0.9.1",
			Date:       time.Now(),
			Version:    "v0.9.1",
			Prerelease: false,
			Assets: []*source.Asset{
				{Name: "test.dmg", Size: 5},
				{Name: "test_32-bit.msi", Size: 5},
				{Name: "test_64-bit.msi", Size: 5},
			},
		},
	}

	for _, release := range want {
		for _, asset := range release.Assets {
			_ = s.UploadAsset(ctx, release.Version, asset.Name, data)
			asset.URL = makeURL(release.Version, asset.Name)
		}
	}

	opt := cmpopts.EquateApproxTime(time.Second)

	t.Run("ListReleases", func(t *testing.T) {
		t.Helper()

		got, err := s.ListReleases(ctx, &source.ListOptions{Prerelease: true})
		if err != nil {
			t.Fatal(err)
		}
		if diff := cmp.Diff(want, got, opt); diff != "" {
			t.Error(diff)
		}
	})

	t.Run("GetRelease", func(t *testing.T) {
		t.Helper()

		got, err := s.GetRelease(ctx, want[0].Version)
		if err != nil {
			t.Fatal(err)
		}
		if diff := cmp.Diff(want[0], got, opt); diff != "" {
			t.Error(diff)
		}
	})

	t.Run("UploadAsset", func(t *testing.T) {
		t.Helper()

		err := s.UploadAsset(ctx, want[0].Version, "test.txt", data)
		if err != nil {
			t.Fatal(err)
		}
	})

	t.Run("DownloadAsset", func(t *testing.T) {
		t.Helper()

		b, err := s.DownloadAsset(ctx, want[0].Version, "test.txt")
		if err != nil {
			t.Fatal(err)
		}
		if !bytes.Equal(data, b) {
			t.Error("should be equal")
		}
	})
}
