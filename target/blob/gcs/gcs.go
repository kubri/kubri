package gcs

import (
	"github.com/abemedia/appcast/target"
	"github.com/abemedia/appcast/target/blob/internal/blob"
	_ "gocloud.dev/blob/gcsblob" // blob driver
)

type Config struct {
	Bucket string
	Folder string
}

func New(c Config) (target.Target, error) {
	return blob.New("gs://"+c.Bucket, c.Folder, "https://storage.googleapis.com/"+c.Bucket)
}
