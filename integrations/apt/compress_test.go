package apt //nolint:testpackage

import (
	"bytes"
	"io"
	"testing"
)

func TestCompress(t *testing.T) {
	want := []byte("test")
	exts := []string{".gz", ".xz", ".lzma", ".zst", ""}

	for _, ext := range exts {
		buf := &bytes.Buffer{}

		w, err := compress(ext)(buf)
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

		r, err := decompress(ext)(buf)
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
