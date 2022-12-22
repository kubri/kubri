package gcs

import (
	"strings"

	"github.com/abemedia/appcast/source"
	"github.com/abemedia/appcast/source/blob/internal/blob"
	_ "gocloud.dev/blob/gcsblob" // blob driver
)

func New(c source.Config) (*source.Source, error) {
	bucket, prefix, _ := strings.Cut(c.Repo, "/")
	return blob.New("gs://"+bucket, prefix, "https://storage.googleapis.com/"+bucket)
}

//nolint:gochecknoinits
func init() { source.Register("gs", New) }
