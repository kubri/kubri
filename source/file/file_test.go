package file_test

import (
	"net/url"
	"os"
	"path/filepath"
	"testing"

	"github.com/kubri/kubri/internal/test"
	"github.com/kubri/kubri/source/file"
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
			path := t.TempDir()

			s, err := file.New(file.Config{Path: path, URL: testCase.url})
			if err != nil {
				t.Fatal(err)
			}

			baseURL := testCase.url
			if baseURL == "" {
				baseURL, _ = url.JoinPath("file:///", filepath.ToSlash(path))
			}

			test.Source(t, s, func(version, asset string) string {
				return baseURL + "/" + version + "/" + asset
			})
		})
	}

	t.Run("New_ResolvePathError", func(t *testing.T) {
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
}
