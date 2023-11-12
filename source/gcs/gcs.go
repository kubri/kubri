package gcs

import (
	"github.com/abemedia/appcast/internal/blob"
	"github.com/abemedia/appcast/source"
	_ "gocloud.dev/blob/gcsblob" // blob driver
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
