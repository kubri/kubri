package azureblob

import (
	_ "gocloud.dev/blob/azureblob" // blob driver

	"github.com/kubri/kubri/internal/blob"
	"github.com/kubri/kubri/target"
)

type Config struct {
	Bucket string
	Folder string
	URL    string
}

func New(c Config) (target.Target, error) {
	return blob.NewTarget("azblob://"+c.Bucket, c.Folder, c.URL)
}
