package azureblob

import (
	"github.com/abemedia/appcast/source"
	"github.com/abemedia/appcast/source/blob/internal/blob"
	_ "gocloud.dev/blob/azureblob" // blob driver
)

type Config struct {
	Bucket string
	Folder string
}

func New(c Config) (*source.Source, error) {
	return blob.New("azblob://"+c.Bucket, c.Folder, "")
}
