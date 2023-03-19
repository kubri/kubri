package blob_test

import (
	"path"
	"testing"

	"github.com/abemedia/appcast/internal/blob"
	"github.com/abemedia/appcast/internal/test"
	_ "gocloud.dev/blob/memblob" // blob driver
)

func TestSource(t *testing.T) {
	s, err := blob.NewSource("mem://", "", "mem:/")
	if err != nil {
		t.Fatal(err)
	}

	test.Source(t, s, func(version, asset string) string {
		return "mem://" + path.Join(version, asset)
	})
}
