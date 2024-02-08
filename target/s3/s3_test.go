package s3_test

import (
	"testing"

	"github.com/kubri/kubri/internal/emulator"
	"github.com/kubri/kubri/internal/test"
	"github.com/kubri/kubri/target/s3"
)

func TestS3(t *testing.T) {
	host := emulator.S3(t, "bucket")

	tgt, err := s3.New(s3.Config{
		Bucket:     "bucket",
		Folder:     "folder",
		Endpoint:   host,
		Region:     "us-east-1",
		DisableSSL: true,
	})
	if err != nil {
		t.Fatal(err)
	}

	test.Target(t, tgt, func(asset string) string {
		return "http://" + host + "/bucket/folder/" + asset
	})
}
