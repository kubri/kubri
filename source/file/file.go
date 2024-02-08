package file

import (
	"net/url"
	"path/filepath"

	_ "gocloud.dev/blob/fileblob" // blob driver

	"github.com/kubri/kubri/internal/blob"
	"github.com/kubri/kubri/source"
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
	url, _ := url.JoinPath("file:///", filepath.ToSlash(path))
	if c.URL == "" {
		c.URL = url
	}
	return blob.NewSource(url, "", c.URL)
}
