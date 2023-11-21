package gcs_test

import (
	"testing"

	"github.com/abemedia/appcast/internal/emulator"
	"github.com/abemedia/appcast/internal/test"
	"github.com/abemedia/appcast/target/gcs"
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
