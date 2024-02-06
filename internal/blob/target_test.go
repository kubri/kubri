package blob_test

import (
	"net/url"
	"testing"

	"github.com/kubri/kubri/internal/blob"
	"github.com/kubri/kubri/internal/test"
)

func TestTarget(t *testing.T) {
	tests := []struct {
		name   string
		prefix string
	}{
		{"Default", ""},
		{"Prefix", "/test/"},
	}
	for _, testCase := range tests {
		t.Run(testCase.name, func(t *testing.T) {
			tgt, err := blob.NewTarget("mem://", testCase.prefix, "http://example.com/downloads")
			if err != nil {
				t.Fatal(err)
			}

			test.Target(t, tgt, func(asset string) string {
				u, _ := url.JoinPath("http://example.com/downloads", testCase.prefix, asset)
				return u
			})
		})
	}
}
