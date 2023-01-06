package azureblob

import (
	"strings"

	"github.com/abemedia/appcast/source"
	"github.com/abemedia/appcast/target"
	"github.com/abemedia/appcast/target/blob/internal/blob"
	_ "gocloud.dev/blob/azureblob" // blob driver
)

func New(c source.Config) (target.Target, error) {
	bucket, prefix, _ := strings.Cut(c.Repo, "/")
	return blob.New("azblob://"+bucket, prefix)
}

//nolint:gochecknoinits
func init() { target.Register("azblob", New) }
