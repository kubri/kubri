package cmd_test

import (
	"io"
	"os"
	"testing"
)

type captured struct {
	f *os.File
}

func (s captured) Bytes() []byte {
	s.f.Seek(0, 0)
	b, _ := io.ReadAll(s.f)
	return b
}

func (s captured) String() string {
	return string(s.Bytes())
}

func (s captured) Reset() {
	s.f.Seek(0, 0)
	s.f.Truncate(0)
}

func capture(t *testing.T, f *os.File) *captured {
	t.Helper()

	file, err := os.CreateTemp(t.TempDir(), "")
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { file.Close() })

	old := *f
	*f = *file
	t.Cleanup(func() { *f = old })

	return &captured{file}
}
