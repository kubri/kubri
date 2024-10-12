// Package gcs provides a target implementation for Google Cloud Storage.
package gcs

import (
	_ "gocloud.dev/blob/gcsblob" // blob driver

	"github.com/kubri/kubri/internal/blob"
	"github.com/kubri/kubri/target"
)

// Config represents the configuration for a Google Cloud Storage target.
type Config struct {
	Bucket string
	Folder string
	URL    string
}

// New returns a new Google Cloud Storage target.
func New(c Config) (target.Target, error) {
	if c.URL == "" {
		c.URL = "https://storage.googleapis.com/" + c.Bucket
	}
	return blob.NewTarget("gs://"+c.Bucket, c.Folder, c.URL)
}
