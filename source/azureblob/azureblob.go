// Package azureblob provides a source implementation for Azure Blob Storage.
package azureblob

import (
	_ "gocloud.dev/blob/azureblob" // blob driver

	"github.com/kubri/kubri/internal/blob"
	"github.com/kubri/kubri/source"
)

// Config represents the configuration for an Azure Blob Storage source.
type Config struct {
	Bucket string
	Folder string
	URL    string
}

// New returns a new Azure Blob Storage source.
func New(c Config) (*source.Source, error) {
	return blob.NewSource("azblob://"+c.Bucket, c.Folder, c.URL)
}
