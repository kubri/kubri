package test

import (
	"io/fs"
	"os"
	"path/filepath"
	"strconv"
	"testing"
)

// update indicates whether golden files should be updated.
var update, _ = strconv.ParseBool(os.Getenv("UPDATE_GOLDEN")) //nolint:gochecknoglobals

type PathFilter func(path string) bool

// Ignore returns a PathFilter that ignores files matching the given globs.
func Ignore(globs ...string) PathFilter {
	return func(path string) bool {
		for _, glob := range globs {
			if ok, _ := filepath.Match(glob, filepath.Base(path)); ok {
				return false
			}
			if ok, _ := filepath.Match(glob, path); ok {
				return false
			}
		}
		return true
	}
}

// Golden compares the result directory to the golden directory and updates the
// golden directory if the UPDATE_GOLDEN environment variable is set.
func Golden(t *testing.T, golden, result string, filter ...PathFilter) {
	t.Helper()

	if !update {
		return
	}

	// Create testdata directory if it doesn't exist.
	if err := os.MkdirAll(golden, 0o750); err != nil {
		t.Fatal("failed to create testdata directory:", err)
	}

	t.Cleanup(func() {
		if err := os.RemoveAll(golden); err != nil {
			t.Fatal("failed to remove golden files:", err)
		}

		err := fs.WalkDir(os.DirFS(result), ".", func(path string, d fs.DirEntry, err error) error {
			if err != nil || d.IsDir() {
				return err
			}

			for _, f := range filter {
				if !f(path) {
					return nil
				}
			}

			if err = os.MkdirAll(filepath.Join(golden, filepath.Dir(path)), 0o750); err != nil {
				return err
			}
			if err = os.Rename(filepath.Join(result, path), filepath.Join(golden, path)); err != nil {
				return err
			}
			return nil
		})
		if err != nil {
			t.Fatal("failed to update golden files:", err)
		}
	})
}

func GoldenFile(t *testing.T, file string, data []byte) {
	t.Helper()

	if !update {
		return
	}

	if err := os.WriteFile(file, data, 0o600); err != nil {
		t.Fatal("failed to write golden file:", err)
	}
}
