package s3

import (
	"net/url"

	"github.com/abemedia/appcast/source"
	"github.com/abemedia/appcast/target"
	"github.com/abemedia/appcast/target/blob/internal/blob"
	_ "gocloud.dev/blob/s3blob" // blob driver
)

func New(c source.Config) (target.Target, error) {
	u, err := url.Parse("s3://" + c.Repo)
	if err != nil {
		return nil, err
	}
	prefix := u.Path
	u.Path = ""
	return blob.New(u.String(), prefix)
}

//nolint:gochecknoinits
func init() { target.Register("s3", New) }
