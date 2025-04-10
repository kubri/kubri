package arch

import (
	"compress/gzip"
	"errors"
	"io"

	"github.com/klauspost/compress/zstd"
	"github.com/ulikunitz/xz"
)

func decompress(ext string) func(io.Reader) (io.ReadCloser, error) {
	switch ext {
	case ".zst":
		return func(r io.Reader) (io.ReadCloser, error) {
			rd, err := zstd.NewReader(r)
			return readCloser{rd}, err
		}
	case ".gz":
		return func(r io.Reader) (io.ReadCloser, error) {
			return gzip.NewReader(r)
		}
	case ".xz":
		return func(r io.Reader) (io.ReadCloser, error) {
			rd, err := xz.NewReader(r)
			return io.NopCloser(rd), err
		}
	default:
		return func(io.Reader) (io.ReadCloser, error) {
			return nil, errors.New("unsupported compression format: " + ext)
		}
	}
}

type rc interface {
	io.Reader
	Close()
}

type readCloser struct{ rc }

func (c readCloser) Close() error {
	c.rc.Close()
	return nil
}
