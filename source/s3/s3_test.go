package s3_test

import (
	"path"
	"testing"

	"github.com/kubri/kubri/internal/emulator"
	"github.com/kubri/kubri/internal/test"
	"github.com/kubri/kubri/source/s3"
)

func TestS3(t *testing.T) {
	endpoint := emulator.S3(t, "bucket")

	s, err := s3.New(s3.Config{
		Bucket:   "bucket",
		Folder:   "folder",
		Region:   "us-east-1",
		Endpoint: endpoint,
	})
	if err != nil {
		t.Fatal(err)
	}

	test.Source(t, s, func(version, asset string) string {
		return endpoint + "/bucket/folder/" + path.Join(version, asset)
	})
}
