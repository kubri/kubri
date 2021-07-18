// Package memory is an in-memory simulator of a source backend for use in tests.
package memory

import (
	"fmt"
	"time"

	"github.com/abemedia/appcast/source"
	"golang.org/x/mod/semver"
)

type memorySource struct {
	data map[string]map[string][]byte
}

func New(source.Config) (*source.Source, error) {
	s := &memorySource{data: map[string]map[string][]byte{}}
	return &source.Source{Provider: s}, nil
}

func (s *memorySource) ListReleases() ([]*source.Release, error) {
	r := make([]*source.Release, 0, len(s.data))
	for version, assets := range s.data {
		r = append(r, s.parseRelease(version, assets))
	}

	return r, nil
}

func (s *memorySource) GetRelease(version string) (*source.Release, error) {
	if assets, ok := s.data[version]; ok {
		return s.parseRelease(version, assets), nil
	}

	return nil, source.ErrReleaseNotFound
}

func (s *memorySource) parseRelease(version string, assets map[string][]byte) *source.Release {
	r := &source.Release{
		Name:       version,
		Version:    version,
		Date:       time.Now(),
		Prerelease: semver.Prerelease(version) != "",
		Assets:     make([]*source.Asset, 0, len(assets)),
	}

	for name, asset := range assets {
		r.Assets = append(r.Assets, &source.Asset{
			Name: name,
			URL:  fmt.Sprintf("memory://%s/%s", version, name),
			Size: len(asset),
		})
	}

	return r
}

func (s *memorySource) UploadAsset(version, name string, data []byte) error {
	r, ok := s.data[version]
	if !ok {
		r = map[string][]byte{}
		s.data[version] = r
	}
	r[name] = data

	return nil
}

func (s *memorySource) DownloadAsset(version, name string) ([]byte, error) {
	r, ok := s.data[version][name]
	if !ok {
		return nil, source.ErrAssetNotFound
	}

	return r, nil
}
