package gcs

import (
	_ "gocloud.dev/blob/gcsblob" // blob driver

	"github.com/kubri/kubri/internal/blob"
	"github.com/kubri/kubri/target"
)

type Config struct {
	Bucket string
	Folder string
	URL    string
}

func New(c Config) (target.Target, error) {
	if c.URL == "" {
		c.URL = "https://storage.googleapis.com/" + c.Bucket
	}
	return blob.NewTarget("gs://"+c.Bucket, c.Folder, c.URL)
}
