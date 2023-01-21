package file

import (
	"context"
	"io"
	"os"
	"path/filepath"

	"github.com/abemedia/appcast/target"
)

type Config struct {
	Path string
}

type fileTarget struct {
	path string
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
	return &fileTarget{path}, nil
}

func (s *fileTarget) NewWriter(ctx context.Context, filename string) (io.WriteCloser, error) {
	path := filepath.Join(s.path, filename)
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return nil, err
	}
	return os.Create(path)
}

func (s *fileTarget) NewReader(ctx context.Context, filename string) (io.ReadCloser, error) {
	return os.Open(filepath.Join(s.path, filename))
}

func (s *fileTarget) Sub(dir string) target.Target {
	return &fileTarget{path: filepath.Join(s.path, dir)}
}

func (s *fileTarget) URL(ctx context.Context, filename string) (string, error) {
	return "file://" + filepath.Join(s.path, filename), nil
}
