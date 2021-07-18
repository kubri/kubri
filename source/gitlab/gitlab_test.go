package gitlab_test

import (
	"testing"
	"time"

	"github.com/abemedia/appcast/source"
	"github.com/abemedia/appcast/source/gitlab"
	"github.com/google/go-cmp/cmp"
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

	r, err := gitlab.New(source.Config{Repo: "abemedia/appcast-test"})
	if err != nil {
		t.Fatal(err)
	}

	got, err := r.Releases()
	if err != nil {
		t.Fatal(err)
	}

	if diff := cmp.Diff(want, got); diff != "" {
		t.Error(diff)
	}
}
