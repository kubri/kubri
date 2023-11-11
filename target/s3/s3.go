package s3

import (
	"net/url"

	"github.com/abemedia/appcast/internal/blob"
	"github.com/abemedia/appcast/target"
	_ "gocloud.dev/blob/s3blob" // blob driver
)

type Config struct {
	Bucket     string
	Folder     string
	Endpoint   string
	Region     string
	DisableSSL bool
	URL        string
}

func New(c Config) (target.Target, error) {
	q := url.Values{}
	if c.Region != "" {
		q.Add("region", c.Region)
	}
	if c.DisableSSL {
		q.Add("disableSSL", "true")
	}
	if c.Endpoint != "" {
		q.Add("endpoint", c.Endpoint)
		q.Add("s3ForcePathStyle", "true")
	}
	return blob.NewTarget("s3://"+c.Bucket+"?"+q.Encode(), c.Folder, c.URL)
}
