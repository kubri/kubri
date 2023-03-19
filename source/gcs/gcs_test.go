package gcs_test

import (
	"path"
	"testing"

	"github.com/abemedia/appcast/internal/emulator"
	"github.com/abemedia/appcast/internal/test"
	"github.com/abemedia/appcast/source/gcs"
)

func TestS3(t *testing.T) {
	emulator.GCS(t, "bucket")

	s, err := gcs.New(gcs.Config{Bucket: "bucket", Folder: "folder"})
	if err != nil {
		t.Fatal(err)
	}

	test.Source(t, s, func(version, asset string) string {
		return "https://storage.googleapis.com/bucket/folder/" + path.Join(version, asset)
	})
}
