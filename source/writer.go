package source

import (
	"bytes"
	"context"
	"io"
)

type Writer struct {
	version  string
	filename string
	s        *Source
	buf      *bytes.Buffer
}

func NewWriter(s *Source, version, filename string) io.WriteCloser {
	return &Writer{
		version:  version,
		filename: filename,
		s:        s,
		buf:      &bytes.Buffer{},
	}
}

func (w *Writer) Write(p []byte) (int, error) {
	return w.buf.Write(p)
}

func (w *Writer) Close() error {
	return w.s.UploadAsset(context.Background(), w.version, w.filename, w.buf.Bytes())
}
