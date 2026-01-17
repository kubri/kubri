// Package local is a source backend for loading packages from a single directory.
package local

import (
	"context"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/kubri/kubri/source"
)

type Config struct {
	Path    string
	Version string
}

type localSource struct {
	path    string
	root    string
	version string
}

func New(c Config) (*source.Source, error) {
	root, err := getRoot(c.Path)
	if err != nil {
		return nil, err
	}
	return source.New(&localSource{path: c.Path, root: root, version: c.Version}), nil
}

func (s *localSource) ListReleases(ctx context.Context) ([]*source.Release, error) {
	r, err := s.GetRelease(ctx, s.version)
	if err != nil {
		return nil, err
	}

	return []*source.Release{r}, nil
}

func (s *localSource) GetRelease(ctx context.Context, version string) (*source.Release, error) {
	files, err := getFiles(s.path)
	if err != nil {
		return nil, err
	}

	r := &source.Release{
		Version: version,
		Date:    time.Now(),
		Assets:  make([]*source.Asset, 0, len(files)),
	}

	for _, path := range files {
		if err := ctx.Err(); err != nil {
			return nil, err
		}

		f, err := os.Stat(path)
		if err != nil {
			return nil, err
		}

		path, err = filepath.Abs(path)
		if err != nil {
			return nil, err
		}

		r.Assets = append(r.Assets, &source.Asset{
			Name: f.Name(),
			URL:  "file://" + path,
			Size: int(f.Size()),
		})
	}

	return r, nil
}

func (s *localSource) UploadAsset(_ context.Context, _, name string, data []byte) error {
	return os.WriteFile(filepath.Join(s.root, name), data, os.ModePerm) //nolint:gosec
}

func (s *localSource) DownloadAsset(_ context.Context, _, name string) ([]byte, error) {
	path := filepath.Join(s.root, name)
	if _, err := os.Stat(path); err == nil {
		return os.ReadFile(path)
	}

	files, err := getFiles(s.path)
	if err != nil {
		return nil, err
	}

	for _, path := range files {
		if filepath.Base(path) == name {
			return os.ReadFile(path)
		}
	}

	return nil, source.ErrAssetNotFound
}

func getFiles(path string) ([]string, error) {
	if strings.ContainsRune(path, '*') {
		return filepath.Glob(path)
	}

	fi, err := os.Stat(path)
	if err != nil {
		return nil, err
	}

	if !fi.IsDir() {
		return []string{path}, nil
	}

	files, err := os.ReadDir(path)
	if err != nil {
		return nil, err
	}

	paths := make([]string, 0, len(files))
	for _, file := range files {
		if !file.IsDir() {
			paths = append(paths, filepath.Join(path, file.Name()))
		}
	}

	return paths, nil
}

func getRoot(path string) (string, error) {
	if i := strings.IndexByte(path, '*'); i >= 0 { //nolint:modernize // more performant
		return path[:i], nil
	}

	fi, err := os.Stat(path)
	if err != nil {
		return "", err
	}

	if !fi.IsDir() {
		return filepath.Dir(path), nil
	}

	return path, nil
}
