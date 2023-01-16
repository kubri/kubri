package gcs_test

import (
	"path"
	"testing"

	"github.com/abemedia/appcast/internal/emulator"
	"github.com/abemedia/appcast/source/blob/gcs"
	"github.com/abemedia/appcast/source/blob/internal/test"
)

func TestS3(t *testing.T) {
	emulator.GCS(t, "bucket")

	s, err := gcs.New(gcs.Config{Bucket: "bucket", Folder: "folder"})
	if err != nil {
		t.Fatal(err)
	}

	test.Run(t, s, func(version, asset string) string {
		return "https://storage.googleapis.com/bucket/folder/" + path.Join(version, asset)
	})
}
