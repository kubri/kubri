package file

import (
	"path/filepath"

	"github.com/abemedia/appcast/source"
	"github.com/abemedia/appcast/source/blob/internal/blob"
	_ "gocloud.dev/blob/fileblob" // blob driver
)

func New(c source.Config) (*source.Source, error) {
	path, err := filepath.Abs(c.Repo)
	if err != nil {
		return nil, err
	}
	url := "file://" + path
	return blob.New(url, "", url)
}

//nolint:gochecknoinits
func init() { source.Register("file", New) }
