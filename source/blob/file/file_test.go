package file_test

import (
	"path/filepath"
	"testing"

	"github.com/abemedia/appcast/source"
	"github.com/abemedia/appcast/source/blob/file"
	"github.com/abemedia/appcast/source/blob/internal/testutils"
)

func TestBlobFile(t *testing.T) {
	dir := t.TempDir()
	s, err := file.New(source.Config{Repo: dir})
	if err != nil {
		t.Fatal(err)
	}

	makeURL := func(version, asset string) string {
		return "file://" + filepath.Join(dir, version, asset)
	}

	testutils.TestBlob(t, s, makeURL)
}
