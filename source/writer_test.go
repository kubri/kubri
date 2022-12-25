package source_test

import (
	"bytes"
	"testing"

	"github.com/abemedia/appcast/source"
	"github.com/abemedia/appcast/source/blob/memory"
)

func TestWriter(t *testing.T) {
	s, _ := memory.New(source.Config{})
	w := source.NewWriter(s, "v1.0.0", "test.txt")

	data := []byte("test")

	n, err := w.Write(data)
	if err != nil {
		t.Fatal(err)
	}

	if n != len(data) {
		t.Error("should be input length")
	}

	if err = w.Close(); err != nil {
		t.Fatal(err)
	}

	b, err := s.DownloadAsset("v1.0.0", "test.txt")
	if err != nil {
		t.Fatal(err)
	}

	if !bytes.Equal(data, b) {
		t.Fatalf("should be '%s': %s", data, b)
	}
}
