package testsource_test

import (
	"path/filepath"
	"testing"
	"time"

	"github.com/abemedia/appcast/internal/test"
	"github.com/abemedia/appcast/internal/testsource"
	"github.com/abemedia/appcast/source"
)

func TestFile(t *testing.T) {
	s := testsource.New([]*source.Release{
		{Name: "v1.0.0", Date: time.Now(), Version: "v1.0.0"},
		{Name: "v1.0.0-alpha1", Date: time.Now(), Version: "v1.0.0-alpha1"},
		{Name: "v0.9.1", Date: time.Now(), Version: "v0.9.1"},
	})

	test.Source(t, s, func(version, asset string) string {
		return "https://example.com/" + filepath.Join(version, asset)
	})
}
