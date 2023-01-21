package memory_test

import (
	"testing"

	"github.com/abemedia/appcast/target/blob/memory"
	"github.com/abemedia/appcast/target/internal/test"
)

func TestMemory(t *testing.T) {
	tgt, err := memory.New(memory.Config{})
	if err != nil {
		t.Fatal(err)
	}

	test.Run(t, tgt, func(asset string) string {
		return "mem://" + asset
	})
}
