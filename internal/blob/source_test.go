package blob_test

import (
	"net/url"
	"testing"

	"github.com/abemedia/appcast/internal/blob"
	"github.com/abemedia/appcast/internal/test"
	_ "gocloud.dev/blob/memblob" // blob driver
)

func TestSource(t *testing.T) {
	tests := []struct {
		name   string
		prefix string
	}{
		{"Default", ""},
		{"Prefix", "/test/"},
	}
	for _, testCase := range tests {
		t.Run(testCase.name, func(t *testing.T) {
			s, err := blob.NewSource("mem://", testCase.prefix, "http://example.com/downloads")
			if err != nil {
				t.Fatal(err)
			}

			test.Source(t, s, func(version, asset string) string {
				u, _ := url.JoinPath("http://example.com/downloads", testCase.prefix, version, asset)
				return u
			})
		})
	}
}
