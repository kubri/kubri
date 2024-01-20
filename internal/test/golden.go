package test

import (
	"io/fs"
	"os"
	"path/filepath"
	"strconv"
	"testing"
)

// Update indicates whether golden files should be updated.
var Update, _ = strconv.ParseBool(os.Getenv("UPDATE_GOLDEN")) //nolint:gochecknoglobals

// Golden compares the result directory to the golden directory and updates the
// golden directory if the UPDATE_GOLDEN environment variable is set.
func Golden(t *testing.T, golden, result string) {
	t.Helper()

	if !Update {
		return
	}

	t.Cleanup(func() {
		from := os.DirFS(result)
		to := os.DirFS(golden)
		err := fs.WalkDir(from, ".", func(path string, d fs.DirEntry, err error) error {
			if err != nil || d.IsDir() {
				return err
			}
			if info, err := fs.Stat(to, path); err == nil {
				b, err := fs.ReadFile(from, path)
				if err != nil {
					return err
				}
				if err = os.WriteFile(filepath.Join(golden, path), b, info.Mode()); err != nil {
					return err
				}
			}
			return nil
		})
		if err != nil {
			t.Log("failed to update golden files:", err)
		}
	})
}
