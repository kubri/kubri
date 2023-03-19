package azureblob

import (
	"github.com/abemedia/appcast/internal/blob"
	"github.com/abemedia/appcast/target"
	_ "gocloud.dev/blob/azureblob" // blob driver
)

type Config struct {
	Bucket string
	Folder string
}

func New(c Config) (target.Target, error) {
	return blob.NewTarget("azblob://"+c.Bucket, c.Folder, "")
}
