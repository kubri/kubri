package file

import (
	"net/url"
	"path/filepath"

	"github.com/abemedia/appcast/internal/blob"
	"github.com/abemedia/appcast/source"
	_ "gocloud.dev/blob/fileblob" // blob driver
)

type Config struct {
	Path string
	URL  string
}

func New(c Config) (*source.Source, error) {
	path, err := filepath.Abs(c.Path)
	if err != nil {
		return nil, err
	}
	url, err := url.JoinPath("file:///", filepath.ToSlash(path))
	if err != nil {
		return nil, err
	}
	if c.URL == "" {
		c.URL = url
	}
	return blob.NewSource(url, "", c.URL)
}
