package file

import (
	"fmt"
	"io/fs"
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

	path, err := filepath.Abs(c.Repo)
	if err != nil {
		return nil, err
	}

	return &source.Source{Provider: &fileSource{path: path}}, nil
}

func (s *fileSource) ListReleases() ([]*source.Release, error) {
	dirs, err := os.ReadDir(s.path)
	if err != nil {
		return nil, err
	}

	result := make([]*source.Release, 0, len(dirs))
	for _, dir := range dirs {
		dirInfo, err := dir.Info()
		if err != nil {
			continue
		}

		r, err := s.getRelease(dirInfo)
		if err != nil {
			continue
		}

		result = append(result, r)
	}

	return result, nil
}

func (s *fileSource) GetRelease(version string) (*source.Release, error) {
	dirInfo, err := os.Stat(filepath.Join(s.path, version))
	if err != nil {
		return nil, err
	}

	return s.getRelease(dirInfo)
}

func (s *fileSource) getRelease(dir fs.FileInfo) (*source.Release, error) {
	if !dir.IsDir() {
		return nil, source.ErrReleaseNotFound
	}

	path := filepath.Join(s.path, dir.Name())
	files, err := os.ReadDir(path)
	if err != nil {
		return nil, err
	}

	r := &source.Release{
		Version: dir.Name(),
		Date:    dir.ModTime(),
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

	return r, nil
}

func (s *fileSource) UploadAsset(version, name string, data []byte) error {
	return os.WriteFile(filepath.Join(s.path, version, name), data, os.ModePerm)
}

func (s *fileSource) DownloadAsset(version, name string) ([]byte, error) {
	return os.ReadFile(filepath.Join(s.path, version, name))
}

//nolint:gochecknoinits
func init() { source.Register("file", New) }
