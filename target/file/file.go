package file

import (
	"context"
	"io"
	"os"
	"path/filepath"

	"github.com/abemedia/appcast/source"
	"github.com/abemedia/appcast/target"
)

type fileTarget struct {
	path string
}

func New(c source.Config) (target.Target, error) {
	err := os.MkdirAll(c.Repo, 0o755)
	if err != nil {
		return nil, err
	}
	path, err := filepath.Abs(c.Repo)
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

//nolint:gochecknoinits
func init() { target.Register("file", New) }
