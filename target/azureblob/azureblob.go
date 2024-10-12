// Package azureblob provides a target implementation for Azure Blob Storage.
package azureblob

import (
	_ "gocloud.dev/blob/azureblob" // blob driver

	"github.com/kubri/kubri/internal/blob"
	"github.com/kubri/kubri/target"
)

// Config represents the configuration for an Azure Blob Storage target.
type Config struct {
	Bucket string
	Folder string
	URL    string
}

// New returns a new Azure Blob Storage target.
func New(c Config) (target.Target, error) {
	return blob.NewTarget("azblob://"+c.Bucket, c.Folder, c.URL)
}
