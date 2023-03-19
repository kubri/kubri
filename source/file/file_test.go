package file_test

import (
	"path/filepath"
	"testing"

	"github.com/abemedia/appcast/internal/test"
	"github.com/abemedia/appcast/source/file"
)

func TestFile(t *testing.T) {
	path := t.TempDir()

	s, err := file.New(file.Config{Path: path})
	if err != nil {
		t.Fatal(err)
	}

	test.Source(t, s, func(version, asset string) string {
		return "file://" + filepath.Join(path, version, asset)
	})
}
