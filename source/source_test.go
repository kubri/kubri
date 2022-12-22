package source_test

import (
	"testing"

	"github.com/abemedia/appcast/source"
	"github.com/google/go-cmp/cmp"
)

type fakeSource struct {
	source.Config
}

func (*fakeSource) ListReleases() ([]*source.Release, error) {
	return []*source.Release{
		{Version: "v0.9.0"},
		{Version: "v1.0.0-pre"},
		{Version: "v1.0.0"},
	}, nil
}

func (*fakeSource) GetRelease(version string) (*source.Release, error) {
	return &source.Release{Version: version}, nil
}

func (*fakeSource) UploadAsset(version, name string, data []byte) error {
	return nil
}

func (*fakeSource) DownloadAsset(version, name string) ([]byte, error) {
	return nil, nil
}

func TestSource(t *testing.T) {
	want := []*source.Release{
		{
			Name:    "v1.0.0",
			Version: "v1.0.0",
		},
		{
			Name:       "v1.0.0-pre",
			Version:    "v1.0.0-pre",
			Prerelease: true,
		},
		{
			Name:    "v0.9.0",
			Version: "v0.9.0",
		},
	}

	r := (&source.Source{&fakeSource{}})

	t.Run("ListReleases", func(t *testing.T) {
		got, _ := r.ListReleases(nil)
		if diff := cmp.Diff(want, got); diff != "" {
			t.Error(diff)
		}
	})

	t.Run("ListReleasesConstraint", func(t *testing.T) {
		got, _ := r.ListReleases(&source.ListOptions{Constraint: "v1"})
		if diff := cmp.Diff(want[:1], got); diff != "" {
			t.Error(diff)
		}
	})

	t.Run("GetRelease", func(t *testing.T) {
		got, _ := r.GetRelease(want[0].Version)
		if diff := cmp.Diff(want[0], got); diff != "" {
			t.Error(diff)
		}
	})
}

func TestSourceUnmarshal(t *testing.T) {
	source.Register("fake", func(c source.Config) (*source.Source, error) {
		return &source.Source{&fakeSource{c}}, nil
	})

	t.Setenv("FAKE_TOKEN", "fake")
	s := &source.Source{}
	if err := s.UnmarshalText([]byte("fake://user/repo")); err != nil {
		t.Fatal(err)
	}

	want := &source.Source{&fakeSource{source.Config{Repo: "user/repo", Token: "fake"}}}
	if diff := cmp.Diff(want, s); diff != "" {
		t.Error(diff)
	}
}
