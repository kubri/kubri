package gcs

import (
	"github.com/abemedia/appcast/internal/blob"
	"github.com/abemedia/appcast/target"
	_ "gocloud.dev/blob/gcsblob" // blob driver
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
