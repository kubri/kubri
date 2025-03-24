package test

import (
	"errors"
	"io"
	"io/fs"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"

	"github.com/kubri/kubri/target"
)

type targetOptions struct {
	delay                time.Duration
	ignoreRemoveNotFound bool
}

type TargetOption func(*targetOptions)

// WithDelay adds a delay between writing, updating & reading, to allow for
// services where the changes aren't available instantly.
func WithDelay(d time.Duration) TargetOption {
	return func(opts *targetOptions) {
		opts.delay = d
	}
}

// WithIgnoreRemoveNotFound disables testing if removing a non-existent file
// returns an error.
func WithIgnoreRemoveNotFound() TargetOption {
	return func(opts *targetOptions) {
		opts.ignoreRemoveNotFound = true
	}
}

// Target tests the given target.
//
//nolint:funlen,gocognit
func Target(t *testing.T, tgt target.Target, makeURL func(string) string, opt ...TargetOption) {
	t.Helper()

	opts := &targetOptions{}
	for _, o := range opt {
		o(opts)
	}

	data := []byte("test")

	t.Run("NewWriter_Create", func(t *testing.T) {
		t.Helper()

		w, err := tgt.NewWriter(t.Context(), "path/to/file")
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

	time.Sleep(opts.delay)

	t.Run("NewWriter_Update", func(t *testing.T) {
		t.Helper()

		w, err := tgt.NewWriter(t.Context(), "path/to/file")
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

	time.Sleep(opts.delay)

	t.Run("NewReader", func(t *testing.T) {
		t.Helper()

		r, err := tgt.NewReader(t.Context(), "path/to/file")
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

		_, err = tgt.NewReader(t.Context(), "does/not/exist")
		if !errors.Is(err, fs.ErrNotExist) {
			t.Fatalf("Should return %q - got %q", fs.ErrNotExist, err)
		}
	})

	t.Run("Sub", func(t *testing.T) {
		t.Helper()

		sub := tgt.Sub("path/to")

		r, err := sub.NewReader(t.Context(), "file")
		if err != nil {
			t.Fatal(err)
		} else {
			_ = r.Close()
		}
	})

	t.Run("URL", func(t *testing.T) {
		t.Helper()

		url, err := tgt.URL(t.Context(), "path/to/file")
		if err != nil {
			t.Fatal(err)
		}

		if diff := cmp.Diff(makeURL("path/to/file"), url); diff != "" {
			t.Fatal(diff)
		}
	})

	t.Run("Remove", func(t *testing.T) {
		t.Helper()

		err := tgt.Remove(t.Context(), "path/to/file")
		if err != nil {
			t.Fatal(err)
		}

		time.Sleep(opts.delay)

		r, err := tgt.NewReader(t.Context(), "path/to/file")
		if err == nil {
			_ = r.Close()
		}
		if !errors.Is(err, fs.ErrNotExist) {
			t.Fatalf("Read should return %q - got %q", fs.ErrNotExist, err)
		}

		if !opts.ignoreRemoveNotFound {
			err = tgt.Remove(t.Context(), "does/not/exist")
			if !errors.Is(err, fs.ErrNotExist) {
				t.Fatalf("Remove should return %q - got %q", fs.ErrNotExist, err)
			}
		}
	})
}
