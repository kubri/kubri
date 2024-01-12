package file_test

import (
	"context"
	"net/url"
	"os"
	"path/filepath"
	"testing"

	"github.com/abemedia/appcast/internal/test"
	"github.com/abemedia/appcast/target/file"
)

func TestFile(t *testing.T) {
	tests := []struct {
		name string
		url  string
	}{
		{"FileURL", ""},
		{"CustomURL", "http://dl.example.com"},
	}

	for _, testCase := range tests {
		t.Run(testCase.name, func(t *testing.T) {
			dir := t.TempDir()

			tgt, err := file.New(file.Config{Path: dir, URL: testCase.url})
			if err != nil {
				t.Fatal(err)
			}

			baseURL := testCase.url
			if baseURL == "" {
				baseURL, _ = url.JoinPath("file:///", filepath.ToSlash(dir))
			}

			test.Target(t, tgt, func(asset string) string {
				return baseURL + "/" + asset
			})
		})
	}

	t.Run("New", func(t *testing.T) {
		t.Run("ResolvePathError", func(t *testing.T) {
			wd, _ := os.Getwd()
			defer os.Chdir(wd)

			dir := t.TempDir()
			os.Chdir(dir)
			os.Remove(dir)

			_, err := file.New(file.Config{})
			if err == nil {
				t.Fatal("should error")
			}
		})

		t.Run("CreateDirError", func(t *testing.T) {
			dir := t.TempDir()
			os.Chmod(dir, 0o000)

			_, err := file.New(file.Config{Path: filepath.Join(dir, "test")})
			if err == nil {
				t.Fatal("should error")
			}
		})
	})

	t.Run("NewWriter_CreateDirError", func(t *testing.T) {
		dir := t.TempDir()
		os.Chmod(dir, 0o000)

		tgt, err := file.New(file.Config{Path: dir})
		if err != nil {
			t.Fatal(err)
		}

		_, err = tgt.NewWriter(context.Background(), filepath.Join(dir, "test", "test"))
		if err == nil {
			t.Fatal("should error")
		}
	})
}
