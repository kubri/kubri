package file

import (
	"path/filepath"

	"github.com/abemedia/appcast/source"
	"github.com/abemedia/appcast/source/blob/internal/blob"
	_ "gocloud.dev/blob/fileblob" // blob driver
)

type Config struct {
	Path string
}

func New(c Config) (*source.Source, error) {
	path, err := filepath.Abs(c.Path)
	if err != nil {
		return nil, err
	}
	url := "file://" + path
	return blob.New(url, "", url)
}
