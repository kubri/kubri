package file

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/abemedia/appcast/source"
)

type fileSource struct {
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

	s := &fileSource{path: c.Repo}

	return &source.Source{Provider: s}, nil
}

func (s *fileSource) Releases() ([]*source.Release, error) {
	dirs, err := os.ReadDir(s.path)
	if err != nil {
		return nil, err
	}

	result := make([]*source.Release, 0, len(dirs))
	for _, dir := range dirs {
		if !dir.IsDir() {
			continue
		}

		dirInfo, err := dir.Info()
		if err != nil {
			return nil, err
		}

		path, err := filepath.Abs(filepath.Join(s.path, dir.Name()))
		if err != nil {
			return nil, err
		}

		files, err := os.ReadDir(path)
		if err != nil {
			return nil, err
		}

		r := &source.Release{
			Version: dir.Name(),
			Date:    dirInfo.ModTime(),
			Assets:  make([]*source.Asset, 0, len(files)),
		}

		for _, file := range files {
			fileInfo, err := file.Info()
			if err != nil {
				return nil, err
			}

			r.Assets = append(r.Assets, &source.Asset{
				Name: file.Name(),
				URL:  "file://" + filepath.Join(path, file.Name()),
				Size: int(fileInfo.Size()),
			})
		}

		result = append(result, r)
	}

	return result, nil
}

func (s *fileSource) UploadAsset(version, name string, data []byte) error {
	return os.WriteFile(filepath.Join(s.path, version, name), data, os.ModePerm)
}

func (s *fileSource) DownloadAsset(version, name string) ([]byte, error) {
	return os.ReadFile(filepath.Join(s.path, version, name))
}

//nolint:gochecknoinits
func init() { source.Register("file", New) }
