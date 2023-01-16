package blob_test

import (
	"path"
	"testing"

	"github.com/abemedia/appcast/source/blob/internal/test"
	"github.com/abemedia/appcast/source/blob/memory"
)

func TestBlob(t *testing.T) {
	s, err := memory.New(memory.Config{})
	if err != nil {
		t.Fatal(err)
	}

	test.Run(t, s, func(version, asset string) string {
		return "mem://" + path.Join(version, asset)
	})
}
