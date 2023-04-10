package test

import (
	"context"
	"io"
	"testing"

	"github.com/abemedia/appcast/target"
	"github.com/google/go-cmp/cmp"
)

//nolint:funlen
func Target(t *testing.T, tgt target.Target, makeURL func(string) string) {
	t.Helper()

	ctx := context.Background()
	data := []byte("test")

	t.Run("NewWriter_Create", func(t *testing.T) {
		t.Helper()

		w, err := tgt.NewWriter(ctx, "path/to/file")
		if err != nil {
			t.Fatal(err)
		}

		if _, err = w.Write([]byte("foo")); err != nil {
			t.Fatal(err)
		}

		if err = w.Close(); err != nil {
			t.Fatal(err)
		}
	})

	t.Run("NewWriter_Update", func(t *testing.T) {
		t.Helper()

		w, err := tgt.NewWriter(ctx, "path/to/file")
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
		t.Helper()

		r, err := tgt.NewReader(ctx, "path/to/file")
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

		if diff := cmp.Diff(data, got); diff != "" {
			t.Fatal(diff)
		}
	})

	t.Run("Sub", func(t *testing.T) {
		t.Helper()

		sub := tgt.Sub("path/to")

		if _, err := sub.NewReader(ctx, "file"); err != nil {
			t.Fatal(err)
		}
	})

	t.Run("URL", func(t *testing.T) {
		t.Helper()

		url, err := tgt.URL(ctx, "path/to/file")
		if err != nil {
			t.Fatal(err)
		}

		if diff := cmp.Diff(makeURL("path/to/file"), url); diff != "" {
			t.Fatal(diff)
		}
	})
}
