package azureblob

import (
	"github.com/abemedia/appcast/internal/blob"
	"github.com/abemedia/appcast/source"
	_ "gocloud.dev/blob/azureblob" // blob driver
)

type Config struct {
	Bucket string
	Folder string
	URL    string
}

func New(c Config) (*source.Source, error) {
	return blob.NewSource("azblob://"+c.Bucket, c.Folder, c.URL)
}
