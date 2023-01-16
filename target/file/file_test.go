package file_test

import (
	"testing"

	"github.com/abemedia/appcast/target/file"
	"github.com/abemedia/appcast/target/internal/test"
)

func TestFile(t *testing.T) {
	tgt, err := file.New(file.Config{Path: t.TempDir()})
	if err != nil {
		t.Fatal(err)
	}

	test.Run(t, tgt)
}
