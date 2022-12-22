package azureblob

import (
	"strings"

	"github.com/abemedia/appcast/source"
	"github.com/abemedia/appcast/source/blob/internal/blob"
	_ "gocloud.dev/blob/azureblob" // blob driver
)

func New(c source.Config) (*source.Source, error) {
	bucket, prefix, _ := strings.Cut(c.Repo, "/")
	return blob.New("azblob://"+bucket, prefix, "")
}

//nolint:gochecknoinits
func init() { source.Register("azblob", New) }
