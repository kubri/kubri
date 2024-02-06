package gcs_test

import (
	"path"
	"testing"

	"github.com/kubri/kubri/internal/emulator"
	"github.com/kubri/kubri/internal/test"
	"github.com/kubri/kubri/source/gcs"
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
