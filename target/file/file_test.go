package file_test

import (
	"path/filepath"
	"testing"

	"github.com/abemedia/appcast/target/file"
	"github.com/abemedia/appcast/target/internal/test"
)

func TestFile(t *testing.T) {
	path := t.TempDir()

	tgt, err := file.New(file.Config{Path: path})
	if err != nil {
		t.Fatal(err)
	}

	test.Run(t, tgt, func(asset string) string {
		return "file://" + filepath.Join(path, asset)
	})
}
