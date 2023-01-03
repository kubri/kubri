package gcs_test

import (
	"testing"

	"github.com/abemedia/appcast/internal/emulator"
	_ "github.com/abemedia/appcast/target/blob/gcs"
	"github.com/abemedia/appcast/target/internal/test"
)

func TestGCS(t *testing.T) {
	emulator.GCS(t, "bucket")
	test.Run(t, "gs://bucket/folder")
}
