package test

import (
	"bytes"
	"context"
	"io"
	"testing"

	"github.com/abemedia/appcast/target"
)

func Run(t *testing.T, url string) {
	t.Helper()

	s, err := target.Open(url)
	if err != nil {
		t.Fatal(err)
	}

	ctx := context.Background()
	data := []byte("test")

	t.Run("NewWriter", func(t *testing.T) {
		w, err := s.NewWriter(ctx, "folder/file")
		if err != nil {
			t.Fatal(err)
		}

		if _, err = w.Write(data); err != nil {
			t.Fatal(err)
		}

		if err = w.Close(); err != nil {
			t.Fatal(err)
		}
	})

	t.Run("NewReader", func(t *testing.T) {
		r, err := s.NewReader(ctx, "folder/file")
		if err != nil {
			t.Fatal(err)
		}

		got, err := io.ReadAll(r)
		if err != nil {
			t.Fatal(err)
		}

		if err = r.Close(); err != nil {
			t.Fatal(err)
		}

		if !bytes.Equal(data, got) {
			t.Fatal("should be equal")
		}
	})

	t.Run("Sub", func(t *testing.T) {
		sub := s.Sub("folder")

		if _, err := sub.NewReader(ctx, "file"); err != nil {
			t.Fatal(err)
		}
	})
}
