package blob_test

import (
	"path/filepath"
	"testing"

	"github.com/abemedia/appcast/source"
	"github.com/abemedia/appcast/source/blob/internal/testutils"
	"github.com/abemedia/appcast/source/blob/memory"
)

func TestBlob(t *testing.T) {
	t.Parallel()

	s, err := memory.New(source.Config{Repo: "test/test"})
	if err != nil {
		t.Fatal(err)
	}

	makeURL := func(version, asset string) string {
		return "mem://" + filepath.Join(version, asset)
	}

	testutils.TestBlob(t, s, makeURL)
}
