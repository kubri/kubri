// Package s3 provides a target implementation for Amazon S3.
package s3

import (
	"net/url"

	_ "gocloud.dev/blob/s3blob" // blob driver

	"github.com/kubri/kubri/internal/blob"
	"github.com/kubri/kubri/target"
)

// Config represents the configuration for an Amazon S3 target.
type Config struct {
	Bucket   string
	Folder   string
	Endpoint string
	Region   string
	URL      string
}

// New returns a new Amazon S3 target.
func New(c Config) (target.Target, error) {
	q := url.Values{}
	if c.Region != "" {
		q.Add("region", c.Region)
	}
	if c.Endpoint != "" {
		q.Add("endpoint", c.Endpoint)
		q.Add("hostname_immutable", "true")
	}
	return blob.NewTarget("s3://"+c.Bucket+"?"+q.Encode(), c.Folder, c.URL)
}
