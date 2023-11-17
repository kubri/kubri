package apt_test

import (
	"bytes"
	"io"
	"testing"

	"github.com/abemedia/appcast/integrations/apt"
)

func TestCompress(t *testing.T) {
	want := []byte("test")
	exts := []string{".gz", ".bz2", ".xz", ".lzma", ".zst", ""}

	for _, ext := range exts {
		buf := &bytes.Buffer{}

		w, err := apt.Compress(ext)(buf)
		if err != nil {
			t.Errorf("compress %s: %s", ext, err)
			continue
		}
		if _, err = w.Write(want); err != nil {
			t.Errorf("write %s: %s", ext, err)
			continue
		}
		if err = w.Close(); err != nil {
			t.Errorf("close %s: %s", ext, err)
			continue
		}

		isCompressed := ext != ""
		if isCompressed == bytes.Equal(want, buf.Bytes()) {
			t.Errorf("%s: should be compressed: %t", ext, isCompressed)
			continue
		}

		r, err := apt.Decompress(ext)(buf)
		if err != nil {
			t.Errorf("decompress %s: %s", ext, err)
			continue
		}
		got, err := io.ReadAll(r)
		if err != nil {
			t.Errorf("read %s: %s", ext, err)
			continue
		}

		if !bytes.Equal(want, got) {
			t.Errorf("%s: should be decompressed", ext)
		}
	}
}
