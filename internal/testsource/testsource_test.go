package testsource_test

import (
	"path/filepath"
	"testing"

	"github.com/abemedia/appcast/internal/test"
	"github.com/abemedia/appcast/internal/testsource"
)

func TestFile(t *testing.T) {
	s := testsource.New(test.SourceWant())
	test.Source(t, s, func(version, asset string) string {
		return "https://example.com/" + filepath.Join(version, asset)
	})
}
