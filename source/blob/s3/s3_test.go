package s3_test

import (
	"path"
	"testing"

	"github.com/abemedia/appcast/internal/emulator"
	"github.com/abemedia/appcast/source/blob/internal/test"
	_ "github.com/abemedia/appcast/source/blob/s3"
)

func TestS3(t *testing.T) {
	host := emulator.S3(t, "bucket")
	repo := "bucket/downloads/test"
	url := "s3://" + repo + "?endpoint=" + host + "&disableSSL=true&s3ForcePathStyle=true"

	test.Run(t, url, func(version, asset string) string {
		return "http://" + host + "/" + path.Join(repo, version, asset)
	})
}
