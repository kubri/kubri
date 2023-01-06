package gcs_test

import (
	"path"
	"testing"

	"github.com/abemedia/appcast/internal/emulator"
	_ "github.com/abemedia/appcast/source/blob/gcs"
	"github.com/abemedia/appcast/source/blob/internal/test"
)

func TestS3(t *testing.T) {
	emulator.GCS(t, "bucket")
	repo := "bucket/downloads/test"

	test.Run(t, "gs://"+repo, func(version, asset string) string {
		return "https://storage.googleapis.com/" + path.Join(repo, version, asset)
	})
}
