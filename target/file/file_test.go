package file_test

import (
	"net/url"
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
}
