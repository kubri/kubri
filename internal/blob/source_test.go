package blob_test

import (
	"bytes"
	"context"
	"path"
	"strings"
	"testing"
	"time"

	"github.com/abemedia/appcast/internal/blob"
	"github.com/abemedia/appcast/internal/test"
	_ "gocloud.dev/blob/memblob" // blob driver
)

func TestSource(t *testing.T) {
	s, err := blob.NewSource("mem://", "", "mem://")
	if err != nil {
		t.Fatal(err)
	}

	w := bytes.NewBufferString("# Changelog")
	for _, v := range test.SourceWant() {
		w.WriteString("\n\n## [" + strings.TrimPrefix(v.Version, "v") + "] - " + v.Date.Format(time.DateOnly))
		if v.Description != "" {
			w.WriteString("\n\n")
			w.WriteString(v.Description)
		}
	}

	for _, v := range test.SourceWant() {
		s.UploadAsset(context.Background(), v.Version, "CHANGELOG.md", w.Bytes())
	}

	test.Source(t, s, func(version, asset string) string {
		return "mem://" + path.Join(version, asset)
	})
}
