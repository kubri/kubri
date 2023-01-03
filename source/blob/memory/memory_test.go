package memory_test

import (
	"path"
	"testing"

	"github.com/abemedia/appcast/source/blob/internal/test"
	_ "github.com/abemedia/appcast/source/blob/memory"
)

func TestMemory(t *testing.T) {
	test.Run(t, "mem://", func(version, asset string) string {
		return "mem://" + path.Join(version, asset)
	})
}
