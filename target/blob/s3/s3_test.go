package s3_test

import (
	"testing"

	"github.com/abemedia/appcast/internal/emulator"
	_ "github.com/abemedia/appcast/target/blob/s3"
	"github.com/abemedia/appcast/target/internal/test"
)

func TestS3(t *testing.T) {
	host := emulator.S3(t, "bucket")
	test.Run(t, "s3://bucket/folder?endpoint="+host+"&disableSSL=true&s3ForcePathStyle=true")
}
