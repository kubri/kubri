package gcs_test

import (
	"testing"

	"github.com/kubri/kubri/internal/emulator"
	"github.com/kubri/kubri/internal/test"
	"github.com/kubri/kubri/target/gcs"
)

func TestGCS(t *testing.T) {
	emulator.GCS(t, "bucket")

	tgt, err := gcs.New(gcs.Config{Bucket: "bucket", Folder: "folder"})
	if err != nil {
		t.Fatal(err)
	}

	test.Target(t, tgt, func(asset string) string {
		return "https://storage.googleapis.com/bucket/folder/" + asset
	}, test.WithIgnoreRemoveNotFound())
}
