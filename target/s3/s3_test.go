package s3_test

import (
	"testing"

	"github.com/kubri/kubri/internal/emulator"
	"github.com/kubri/kubri/internal/test"
	"github.com/kubri/kubri/target/s3"
)

func TestS3(t *testing.T) {
	endpoint := emulator.S3(t, "bucket")

	tgt, err := s3.New(s3.Config{
		Bucket:   "bucket",
		Folder:   "folder",
		Endpoint: endpoint,
		Region:   "us-east-1",
	})
	if err != nil {
		t.Fatal(err)
	}

	test.Target(t, tgt, func(asset string) string {
		return endpoint + "/bucket/folder/" + asset
	})
}
