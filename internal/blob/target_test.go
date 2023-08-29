package blob_test

import (
	"testing"

	"github.com/abemedia/appcast/internal/blob"
	"github.com/abemedia/appcast/internal/test"
)

func TestTarget(t *testing.T) {
	tgt, err := blob.NewTarget("mem://", "", "mem://")
	if err != nil {
		t.Fatal(err)
	}

	test.Target(t, tgt, func(asset string) string {
		return "mem://" + asset
	})
}
