package file_test

import (
	"testing"

	_ "github.com/abemedia/appcast/target/file"
	"github.com/abemedia/appcast/target/internal/test"
)

func TestFile(t *testing.T) {
	test.Run(t, "file://"+t.TempDir())
}
