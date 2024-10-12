// Package s3 provides a source implementation for Amazon S3.
package s3

import (
	"net/url"

	_ "gocloud.dev/blob/s3blob" // blob driver

	"github.com/kubri/kubri/internal/blob"
	"github.com/kubri/kubri/source"
)

// Config represents the configuration for an Amazon S3 source.
type Config struct {
	Bucket   string
	Folder   string
	Endpoint string
	Region   string
	URL      string
}

// New returns a new Amazon S3 source.
func New(c Config) (*source.Source, error) {
	q := url.Values{}
	if c.Region != "" {
		q.Add("region", c.Region)
	}
	if c.Endpoint != "" {
		q.Add("endpoint", c.Endpoint)
		q.Add("hostname_immutable", "true")
	}
	return blob.NewSource("s3://"+c.Bucket+"?"+q.Encode(), c.Folder, c.URL)
}
