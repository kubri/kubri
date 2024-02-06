// Package testsource is an in-memory simulator of a source for use in tests.
package testsource

import (
	"context"

	"github.com/kubri/kubri/source"
)

type testSource struct {
	data   []*source.Release
	assets map[[2]string][]byte
}

func New(r []*source.Release) *source.Source {
	return source.New(&testSource{r, map[[2]string][]byte{}})
}

func (s *testSource) GetRelease(_ context.Context, version string) (*source.Release, error) {
	for _, r := range s.data {
		if r.Version == version {
			return r, nil
		}
	}
	return nil, source.ErrNoReleaseFound
}

func (s *testSource) ListReleases(context.Context) ([]*source.Release, error) {
	return append([]*source.Release(nil), s.data...), nil
}

func (s *testSource) DownloadAsset(ctx context.Context, version, name string) ([]byte, error) {
	r, err := s.GetRelease(ctx, version)
	if err != nil {
		return nil, err
	}

	if a := getAsset(r.Assets, name); a != nil {
		if b := s.assets[[2]string{version, name}]; b != nil {
			return b, nil
		}
	}

	return nil, source.ErrAssetNotFound
}

func (s *testSource) UploadAsset(ctx context.Context, version, name string, data []byte) error {
	r, err := s.GetRelease(ctx, version)
	if err != nil {
		return err
	}

	asset := &source.Asset{
		Name: name,
		Size: len(data),
		URL:  "https://example.com/" + version + "/" + name,
	}

	if a := getAsset(r.Assets, name); a != nil {
		*a = *asset
	} else {
		r.Assets = append(r.Assets, asset)
	}

	s.assets[[2]string{version, name}] = data

	return nil
}

func getAsset(assets []*source.Asset, name string) *source.Asset {
	for _, a := range assets {
		if a.Name == name {
			return a
		}
	}
	return nil
}
