package gcs

import (
	"github.com/abemedia/appcast/source"
	"github.com/abemedia/appcast/source/blob/internal/blob"
	_ "gocloud.dev/blob/gcsblob" // blob driver
)

type Config struct {
	Bucket string
	Folder string
}

func New(c Config) (*source.Source, error) {
	return blob.New("gs://"+c.Bucket, c.Folder, "https://storage.googleapis.com/"+c.Bucket)
}
