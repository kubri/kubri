package s3_test

import (
	"path"
	"testing"

	"github.com/kubri/kubri/internal/emulator"
	"github.com/kubri/kubri/internal/test"
	"github.com/kubri/kubri/source/s3"
)

func TestS3(t *testing.T) {
	host := emulator.S3(t, "bucket")

	s, err := s3.New(s3.Config{
		Bucket:     "bucket",
		Folder:     "folder",
		Region:     "us-east-1",
		Endpoint:   host,
		DisableSSL: true,
	})
	if err != nil {
		t.Fatal(err)
	}

	test.Source(t, s, func(version, asset string) string {
		return "http://" + host + "/bucket/folder/" + path.Join(version, asset)
	})
}
