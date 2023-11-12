package file

import (
	"context"
	"io"
	"net/url"
	"os"
	"path/filepath"

	"github.com/abemedia/appcast/target"
)

type Config struct {
	Path string
	URL  string
}

type fileTarget struct {
	path string
	url  string
}

func New(c Config) (target.Target, error) {
	err := os.MkdirAll(c.Path, 0o755)
	if err != nil {
		return nil, err
	}
	path, err := filepath.Abs(c.Path)
	if err != nil {
		return nil, err
	}
	if c.URL == "" {
		c.URL, _ = url.JoinPath("file:///", filepath.ToSlash(path))
	}
	return &fileTarget{path, c.URL}, nil
}

func (t *fileTarget) NewWriter(_ context.Context, filename string) (io.WriteCloser, error) {
	path := filepath.Join(t.path, filename)
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return nil, err
	}
	return os.Create(path)
}

func (t *fileTarget) NewReader(_ context.Context, filename string) (io.ReadCloser, error) {
	return os.Open(filepath.Join(t.path, filename))
}

func (t *fileTarget) Sub(dir string) target.Target {
	u, _ := url.JoinPath(t.url, dir)
	return &fileTarget{path: filepath.Join(t.path, dir), url: u}
}

func (t *fileTarget) URL(_ context.Context, filename string) (string, error) {
	return url.JoinPath(t.url, filename)
}
