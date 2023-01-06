package memory_test

import (
	"testing"

	_ "github.com/abemedia/appcast/target/blob/memory"
	"github.com/abemedia/appcast/target/internal/test"
)

func TestMemory(t *testing.T) {
	test.Run(t, "mem://")
}
