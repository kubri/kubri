package github_test

import (
	"bytes"
	"context"
	"os"
	"testing"
	"time"

	"github.com/abemedia/appcast/source"
	"github.com/abemedia/appcast/source/github"
	"github.com/google/go-cmp/cmp"
	gh "github.com/google/go-github/github"
	"golang.org/x/oauth2"
)

func TestGithub(t *testing.T) {
	want := []*source.Release{
		{
			Name:        "v1.0.0",
			Description: "This is a stable release.",
			Date:        time.Date(2022, 11, 29, 22, 10, 53, 0, time.UTC),
			Version:     "v1.0.0",
			Assets: []*source.Asset{
				{
					Name: "test.dmg",
					URL:  "https://github.com/abemedia/appcast-test/releases/download/v1.0.0/test.dmg",
					Size: 5,
				},
				{
					Name: "test_32-bit.msi",
					URL:  "https://github.com/abemedia/appcast-test/releases/download/v1.0.0/test_32-bit.msi",
					Size: 5,
				},
				{
					Name: "test_64-bit.msi",
					URL:  "https://github.com/abemedia/appcast-test/releases/download/v1.0.0/test_64-bit.msi",
					Size: 5,
				},
			},
		},
		{
			Name:        "v1.0.0-alpha1",
			Description: "This is a pre-release.",
			Date:        time.Date(2022, 11, 29, 22, 10, 19, 0, time.UTC),
			Version:     "v1.0.0-alpha1",
			Prerelease:  true,
			Assets: []*source.Asset{
				{
					Name: "test.dmg",
					URL:  "https://github.com/abemedia/appcast-test/releases/download/v1.0.0-alpha1/test.dmg",
					Size: 5,
				},
				{
					Name: "test_32-bit.msi",
					URL:  "https://github.com/abemedia/appcast-test/releases/download/v1.0.0-alpha1/test_32-bit.msi",
					Size: 5,
				},
				{
					Name: "test_64-bit.msi",
					URL:  "https://github.com/abemedia/appcast-test/releases/download/v1.0.0-alpha1/test_64-bit.msi",
					Size: 5,
				},
			},
		},
	}

	s, err := github.New(source.Config{Repo: "abemedia/appcast-test", Token: os.Getenv("GITHUB_TOKEN")})
	if err != nil {
		t.Fatal(err)
	}

	t.Run("ListReleases", func(t *testing.T) {
		got, err := s.ListReleases(nil)
		if err != nil {
			t.Fatal(err)
		}

		if diff := cmp.Diff(want, got); diff != "" {
			t.Error(diff)
		}
	})

	t.Run("GetRelease", func(t *testing.T) {
		got, err := s.GetRelease(want[0].Version)
		if err != nil {
			t.Fatal(err)
		}

		if diff := cmp.Diff(want[0], got); diff != "" {
			t.Error(diff)
		}
	})

	data := []byte("test")

	t.Run("UploadAsset", func(t *testing.T) {
		if _, ok := os.LookupEnv("GITHUB_TOKEN"); !ok {
			t.Skip("missing GITHUB_TOKEN")
		}
		err := s.UploadAsset(want[0].Version, "test.txt", data)
		if err != nil {
			t.Fatal(err)
		}
	})

	t.Run("DownloadAsset", func(t *testing.T) {
		if _, ok := os.LookupEnv("GITHUB_TOKEN"); !ok {
			t.Skip("missing GITHUB_TOKEN")
		}
		b, err := s.DownloadAsset(want[0].Version, "test.txt")
		if err != nil {
			t.Fatal(err)
		}

		if !bytes.Equal(data, b) {
			t.Error("should be equal")
		}
	})

	t.Cleanup(func() {
		if _, ok := os.LookupEnv("GITHUB_TOKEN"); !ok {
			return
		}

		ctx := context.Background()
		ts := oauth2.StaticTokenSource(&oauth2.Token{AccessToken: os.Getenv("GITHUB_TOKEN")})
		client := gh.NewClient(oauth2.NewClient(ctx, ts))
		release, _, _ := client.Repositories.GetReleaseByTag(ctx, "abemedia", "appcast-test", want[0].Version)

		for _, asset := range release.Assets {
			if asset.GetName() == "test.txt" {
				_, _ = client.Repositories.DeleteReleaseAsset(ctx, "abemedia", "appcast-test", asset.GetID())
			}
		}
	})
}
