// Package testsource is an in-memory simulator of a source for use in tests.
package testsource

import (
	"context"
	"os"
	"path/filepath"
	"testing"
	"time"
	"unsafe"

	"github.com/abemedia/appcast/source"
	"github.com/abemedia/appcast/source/file"
)

//nolint:gochecknoglobals
var (
	defaultVersions = []string{"v1.0.0", "v1.1.0-beta", "v1.1.0", "v2.0.0"}
	defaultURL      = "https://example.com"
)

type testSource struct {
	releases map[string]*source.Release
	driver   source.Driver
}

type Option func(*options)

func WithReleases(r []*source.Release) Option {
	return func(o *options) {
		o.releases = r
	}
}

func WithVersions(v ...string) Option {
	return func(o *options) {
		for _, v := range v {
			o.releases = append(o.releases, &source.Release{
				Date:    time.Now(),
				Version: v,
			})
		}
	}
}

func WithGenerateAssets(g Generator) Option {
	return func(o *options) {
		o.generators = append(o.generators, g)
	}
}

func WithAssets(name ...string) Option {
	return WithGenerateAssets(func(dir, version string) error {
		for _, name := range name {
			err := os.WriteFile(filepath.Join(dir, name), []byte("test"), 0o600)
			if err != nil {
				return err
			}
		}
		return nil
	})
}

type options struct {
	generators []Generator
	releases   []*source.Release
}

func New(t *testing.T, opts ...Option) *source.Source {
	t.Helper()

	dir := t.TempDir()

	o := &options{}
	for _, opt := range opts {
		opt(o)
	}
	if o.releases == nil {
		WithVersions(defaultVersions...)(o)
	}

	releases := make(map[string]*source.Release, len(o.releases))
	for _, r := range o.releases {
		releases[r.Version] = r
		path := filepath.Join(dir, r.Version)
		if err := os.MkdirAll(path, 0o755); err != nil {
			t.Fatal(err)
		}
		for _, g := range o.generators {
			if err := g(path, r.Version); err != nil {
				t.Fatal(err)
			}
		}
	}

	src, err := file.New(file.Config{Path: dir, URL: defaultURL})
	if err != nil {
		t.Fatal(err)
	}

	drv := (*struct{ s source.Driver })(unsafe.Pointer(src)).s

	return source.New(&testSource{
		releases: releases,
		driver:   drv,
	})
}

func (s *testSource) GetRelease(ctx context.Context, version string) (*source.Release, error) {
	r, ok := s.releases[version]
	if !ok {
		return nil, source.ErrNoReleaseFound
	}

	if rel, err := s.driver.GetRelease(ctx, version); err == nil {
		r.Assets = rel.Assets
	}

	return r, nil
}

func (s *testSource) ListReleases(ctx context.Context) ([]*source.Release, error) {
	r := make([]*source.Release, 0, len(s.releases))
	for v, rel := range s.releases {
		if rr, err := s.GetRelease(ctx, v); err == nil {
			rel.Assets = rr.Assets
		}
		r = append(r, rel)
	}
	return r, nil
}

func (s *testSource) DownloadAsset(ctx context.Context, version, name string) ([]byte, error) {
	return s.driver.DownloadAsset(ctx, version, name)
}

func (s *testSource) UploadAsset(ctx context.Context, version, name string, data []byte) error {
	return s.driver.UploadAsset(ctx, version, name, data)
}
