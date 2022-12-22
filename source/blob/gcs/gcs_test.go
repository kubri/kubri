package gcs_test

import (
	"log"
	"path/filepath"
	"testing"

	"github.com/abemedia/appcast/source"
	"github.com/abemedia/appcast/source/blob/gcs"
	"github.com/abemedia/appcast/source/blob/internal/testutils"
	"github.com/fullstorydev/emulators/storage/gcsemu"
)

func TestBlobS3(t *testing.T) {
	emu, err := gcsemu.NewServer(":0", gcsemu.Options{})
	if err != nil {
		log.Fatal(err)
	}
	emu.InitBucket("bucket")
	t.Setenv("STORAGE_EMULATOR_HOST", emu.Addr)

	repo := "bucket/downloads/test"

	s, err := gcs.New(source.Config{Repo: repo})
	if err != nil {
		t.Fatal(err)
	}

	makeURL := func(version, asset string) string {
		return "https://storage.googleapis.com/" + filepath.Join(repo, version, asset)
	}

	testutils.TestBlob(t, s, makeURL)
}
