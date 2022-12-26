package gitlab_test

import (
	"bytes"
	"os"
	"testing"
	"time"

	"github.com/abemedia/appcast/source"
	"github.com/abemedia/appcast/source/gitlab"
	"github.com/google/go-cmp/cmp"
	gl "github.com/xanzy/go-gitlab"
)

func TestGitlab(t *testing.T) {
	want := []*source.Release{
		{
			Name:        "v1.0.0",
			Description: "This is a stable release.",
			Date:        time.Date(2022, 12, 15, 18, 10, 26, 654000000, time.UTC),
			Version:     "v1.0.0",
			Assets: []*source.Asset{
				{
					Name: "test_64-bit.msi",
					URL:  "https://github.com/abemedia/appcast-test/releases/download/v1.0.0/test_64-bit.msi",
					Size: 5,
				},
				{
					Name: "test_32-bit.msi",
					URL:  "https://github.com/abemedia/appcast-test/releases/download/v1.0.0/test_32-bit.msi",
					Size: 5,
				},
				{
					Name: "test.dmg",
					URL:  "https://github.com/abemedia/appcast-test/releases/download/v1.0.0/test.dmg",
					Size: 5,
				},
			},
		},
		{
			Name:        "v1.0.0-alpha1",
			Description: "This is a pre-release.",
			Date:        time.Date(2022, 12, 15, 18, 9, 32, 340000000, time.UTC),
			Version:     "v1.0.0-alpha1",
			Prerelease:  true,
			Assets: []*source.Asset{
				{
					Name: "test_64-bit.msi",
					URL:  "https://github.com/abemedia/appcast-test/releases/download/v1.0.0-alpha1/test_64-bit.msi",
					Size: 5,
				},
				{
					Name: "test_32-bit.msi",
					URL:  "https://github.com/abemedia/appcast-test/releases/download/v1.0.0-alpha1/test_32-bit.msi",
					Size: 5,
				},
				{
					Name: "test.dmg",
					URL:  "https://github.com/abemedia/appcast-test/releases/download/v1.0.0-alpha1/test.dmg",
					Size: 5,
				},
			},
		},
	}

	s, err := gitlab.New(source.Config{Repo: "abemedia/appcast-test", Token: os.Getenv("GITLAB_TOKEN")})
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
		if _, ok := os.LookupEnv("GITLAB_TOKEN"); !ok {
			t.Skip("missing GITLAB_TOKEN")
		}
		err := s.UploadAsset(want[0].Version, "test.txt", data)
		if err != nil {
			t.Fatal(err)
		}
	})

	t.Run("DownloadAsset", func(t *testing.T) {
		if _, ok := os.LookupEnv("GITLAB_TOKEN"); !ok {
			t.Skip("missing GITLAB_TOKEN")
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
		if _, ok := os.LookupEnv("GITLAB_TOKEN"); !ok {
			return
		}

		client, _ := gl.NewClient(os.Getenv("GITLAB_TOKEN"))
		links, _, _ := client.ReleaseLinks.ListReleaseLinks("abemedia/appcast-test", want[0].Version, nil)

		for _, link := range links {
			if link.Name == "test.txt" {
				_, _, _ = client.ReleaseLinks.DeleteReleaseLink("abemedia/appcast-test", want[0].Version, link.ID)
			}
		}
	})
}
