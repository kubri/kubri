package gcs

import (
	"github.com/abemedia/appcast/internal/blob"
	"github.com/abemedia/appcast/source"
	_ "gocloud.dev/blob/gcsblob" // blob driver
)

type Config struct {
	Bucket string
	Folder string
}

func New(c Config) (*source.Source, error) {
	return blob.NewSource("gs://"+c.Bucket, c.Folder, "https://storage.googleapis.com/"+c.Bucket)
}
