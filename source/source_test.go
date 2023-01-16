package source_test

import (
	"context"
	"testing"

	"github.com/abemedia/appcast/source"
	"github.com/google/go-cmp/cmp"
)

type fakeSource struct{}

func (*fakeSource) ListReleases(ctx context.Context) ([]*source.Release, error) {
	return []*source.Release{
		{Version: "foo"},
		{Version: "v0.9.0"},
		{Version: "v1.0.0-pre"},
		{Version: "v1.0.0"},
	}, nil
}

func (*fakeSource) GetRelease(ctx context.Context, version string) (*source.Release, error) {
	return &source.Release{Version: version}, nil
}

func (*fakeSource) UploadAsset(ctx context.Context, version, name string, data []byte) error {
	return nil
}

func (*fakeSource) DownloadAsset(ctx context.Context, version, name string) ([]byte, error) {
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

	s := source.New(&fakeSource{})
	ctx := context.Background()

	t.Run("ListReleases", func(t *testing.T) {
		tests := []struct {
			msg  string
			want []*source.Release
			opt  *source.ListOptions
		}{
			{"nil options", []*source.Release{want[0], want[2]}, nil},
			{"zero options", []*source.Release{want[0], want[2]}, &source.ListOptions{}},
			{"version >= v1-a", want[:1], &source.ListOptions{Version: ">= v1-a"}},
			{"version >= v1-a; prerelease = true", want[:2], &source.ListOptions{Version: ">= v1-a", Prerelease: true}},
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
}
