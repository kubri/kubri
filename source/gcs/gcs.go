package gcs

import (
	_ "gocloud.dev/blob/gcsblob" // blob driver

	"github.com/kubri/kubri/internal/blob"
	"github.com/kubri/kubri/source"
)

type Config struct {
	Bucket string
	Folder string
	URL    string
}

func New(c Config) (*source.Source, error) {
	if c.URL == "" {
		c.URL = "https://storage.googleapis.com/" + c.Bucket
	}
	return blob.NewSource("gs://"+c.Bucket, c.Folder, c.URL)
}
