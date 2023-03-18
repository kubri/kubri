package source_test

import (
	"context"
	"testing"

	"github.com/abemedia/appcast/source"
	"github.com/google/go-cmp/cmp"
)

type fakeSource struct {
	assets map[[2]string][]byte
}

func (*fakeSource) ListReleases(context.Context) ([]*source.Release, error) {
	return []*source.Release{
		{Version: "foo"},
		{Version: "v0.9.0"},
		{Version: "v1.1.0-pre"},
		{Version: "v1.0.0"},
	}, nil
}

func (*fakeSource) GetRelease(_ context.Context, version string) (*source.Release, error) {
	return &source.Release{Version: version}, nil
}

func (s *fakeSource) UploadAsset(_ context.Context, version, name string, data []byte) error {
	s.assets[[2]string{version, name}] = data
	return nil
}

func (s *fakeSource) DownloadAsset(_ context.Context, version, name string) ([]byte, error) {
	if b, ok := s.assets[[2]string{version, name}]; ok {
		return b, nil
	}
	return nil, source.ErrAssetNotFound
}

func TestSource(t *testing.T) {
	want := []*source.Release{
		{
			Name:       "v1.1.0-pre",
			Version:    "v1.1.0-pre",
			Prerelease: true,
		},
		{
			Name:    "v1.0.0",
			Version: "v1.0.0",
		},
		{
			Name:    "v0.9.0",
			Version: "v0.9.0",
		},
	}

	s := source.New(&fakeSource{map[[2]string][]byte{}})
	ctx := context.Background()

	t.Run("ListReleases", func(t *testing.T) {
		tests := []struct {
			msg  string
			want []*source.Release
			opt  *source.ListOptions
		}{
			{"nil options", want[1:], nil},
			{"zero options", want[1:], &source.ListOptions{}},
			{"version >= v1", want[1:2], &source.ListOptions{Version: ">= v1"}},
			{"version >= v1; prerelease = true", want[:2], &source.ListOptions{Version: ">= v1", Prerelease: true}},
		}

		for _, test := range tests {
			got, _ := s.ListReleases(ctx, test.opt)
			if diff := cmp.Diff(test.want, got); diff != "" {
				t.Errorf("%s\n%s", test.msg, diff)
			}
		}
	})

	t.Run("GetRelease", func(t *testing.T) {
		got, _ := s.GetRelease(ctx, want[0].Version)
		if diff := cmp.Diff(want[0], got); diff != "" {
			t.Error(diff)
		}
	})

	t.Run("UploadAsset", func(t *testing.T) {
		err := s.UploadAsset(ctx, want[0].Version, "test", []byte("test"))
		if err != nil {
			t.Error(err)
		}
	})

	t.Run("DownloadAsset", func(t *testing.T) {
		got, err := s.DownloadAsset(ctx, want[0].Version, "test")
		if err != nil {
			t.Error(err)
		}
		if diff := cmp.Diff("test", string(got)); diff != "" {
			t.Error(diff)
		}
	})
}
