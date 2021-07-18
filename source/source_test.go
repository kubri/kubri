package source_test

import (
	"testing"

	"github.com/abemedia/appcast/source"
	"github.com/google/go-cmp/cmp"
)

type fakeSource struct {
	source.Config
}

func (*fakeSource) Releases() ([]*source.Release, error) {
	return []*source.Release{
		{Version: "v1.0.0"},
		{Version: "v1.0.0-pre"},
	}, nil
}

func (*fakeSource) UploadAsset(version, name string, data []byte) error {
	return nil
}

func (*fakeSource) DownloadAsset(version, name string) ([]byte, error) {
	return nil, nil
}

func TestSourceReleases(t *testing.T) {
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
	}

	got, err := (&source.Source{&fakeSource{}}).Releases()
	if err != nil {
		t.Fatal(err)
	}

	if diff := cmp.Diff(want, got); diff != "" {
		t.Error(diff)
	}
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
