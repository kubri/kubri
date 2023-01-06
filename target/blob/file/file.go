package file

import (
	"os"
	"path/filepath"

	"github.com/abemedia/appcast/source"
	"github.com/abemedia/appcast/target"
	"github.com/abemedia/appcast/target/blob/internal/blob"
	_ "gocloud.dev/blob/fileblob" // blob driver
)

func New(c source.Config) (target.Target, error) {
	err := os.MkdirAll(c.Repo, 0o755)
	if err != nil {
		return nil, err
	}
	path, err := filepath.Abs(c.Repo)
	if err != nil {
		return nil, err
	}
	url := "file://" + path
	return blob.New(url, "")
}

//nolint:gochecknoinits
func init() { target.Register("file", New) }
