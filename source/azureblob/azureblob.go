package azureblob

import (
	_ "gocloud.dev/blob/azureblob" // blob driver

	"github.com/kubri/kubri/internal/blob"
	"github.com/kubri/kubri/source"
)

type Config struct {
	Bucket string
	Folder string
	URL    string
}

func New(c Config) (*source.Source, error) {
	return blob.NewSource("azblob://"+c.Bucket, c.Folder, c.URL)
}
