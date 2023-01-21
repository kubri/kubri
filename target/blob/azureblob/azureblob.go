package azureblob

import (
	"github.com/abemedia/appcast/target"
	"github.com/abemedia/appcast/target/blob/internal/blob"
	_ "gocloud.dev/blob/azureblob" // blob driver
)

type Config struct {
	Bucket string
	Folder string
}

func New(c Config) (target.Target, error) {
	return blob.New("azblob://"+c.Bucket, c.Folder, "")
}
