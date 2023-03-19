package gcs

import (
	"github.com/abemedia/appcast/internal/blob"
	"github.com/abemedia/appcast/target"
	_ "gocloud.dev/blob/gcsblob" // blob driver
)

type Config struct {
	Bucket string
	Folder string
}

func New(c Config) (target.Target, error) {
	return blob.NewTarget("gs://"+c.Bucket, c.Folder, "https://storage.googleapis.com/"+c.Bucket)
}
