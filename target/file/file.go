// Package file provides a target implementation for the local filesystem.
package file

import (
	"context"
	"io"
	"net/url"
	"os"
	"path/filepath"

	"github.com/kubri/kubri/target"
)

// Config represents the configuration for a file target.
type Config struct {
	Path string
	URL  string
}

// New returns a new file target.
func New(c Config) (target.Target, error) {
	path, err := filepath.Abs(c.Path)
	if err != nil {
		return nil, err
	}
	if err = os.MkdirAll(path, 0o750); err != nil {
		return nil, err
	}
	if c.URL == "" {
		c.URL, _ = url.JoinPath("file:///", filepath.ToSlash(path))
	}
	return &fileTarget{path, c.URL}, nil
}

type fileTarget struct {
	path string
	url  string
}

func (t *fileTarget) NewWriter(_ context.Context, filename string) (io.WriteCloser, error) {
	path := filepath.Join(t.path, filename)
	if err := os.MkdirAll(filepath.Dir(path), 0o750); err != nil {
		return nil, err
	}
	return os.Create(path)
}

func (t *fileTarget) NewReader(_ context.Context, filename string) (io.ReadCloser, error) {
	return os.Open(filepath.Join(t.path, filename))
}

func (t *fileTarget) Remove(_ context.Context, filename string) error {
	return os.Remove(filepath.Join(t.path, filename))
}

func (t *fileTarget) Sub(dir string) target.Target {
	u, _ := url.JoinPath(t.url, dir)
	return &fileTarget{path: filepath.Join(t.path, dir), url: u}
}

func (t *fileTarget) URL(_ context.Context, filename string) (string, error) {
	return url.JoinPath(t.url, filename)
}
