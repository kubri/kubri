package local

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/abemedia/appcast/source"
)

type localSource struct {
	path string
}

func New(c source.Config) (*source.Source, error) {
	fs, err := os.Stat(c.Repo)
	if err != nil {
		return nil, err
	}

	if !fs.IsDir() {
		return nil, fmt.Errorf("not a directory: %s", c.Repo)
	}

	s := &localSource{path: c.Repo}

	return &source.Source{Provider: s}, nil
}

func (s *localSource) Releases() ([]*source.Release, error) {
	files, err := os.ReadDir(s.path)
	if err != nil {
		return nil, err
	}

	r := &source.Release{
		Version: "v0.0.0",
		Date:    time.Now(),
		Assets:  make([]*source.Asset, 0, len(files)),
	}

	for _, file := range files {
		fileInfo, err := file.Info()
		if err != nil {
			return nil, err
		}

		r.Assets = append(r.Assets, &source.Asset{
			Name: file.Name(),
			URL:  "file://" + filepath.Join(s.path, file.Name()),
			Size: int(fileInfo.Size()),
		})
	}

	return []*source.Release{r}, nil
}

func (s *localSource) UploadAsset(_, name string, data []byte) error {
	return os.WriteFile(filepath.Join(s.path, name), data, os.ModePerm)
}

func (s *localSource) DownloadAsset(_, name string) ([]byte, error) {
	return os.ReadFile(filepath.Join(s.path, name))
}

//nolint:gochecknoinits
func init() { source.Register("local", New) }
