package file_test

import (
	"path/filepath"
	"testing"

	_ "github.com/abemedia/appcast/source/blob/file"
	"github.com/abemedia/appcast/source/blob/internal/test"
)

func TestFile(t *testing.T) {
	dir := t.TempDir()

	test.Run(t, "file://"+dir, func(version, asset string) string {
		return "file://" + filepath.Join(dir, version, asset)
	})
}
