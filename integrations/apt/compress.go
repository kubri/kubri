package apt

import (
	"compress/bzip2"
	"compress/gzip"
	"io"

	"github.com/klauspost/compress/zstd"
	"github.com/ulikunitz/xz"
	"github.com/ulikunitz/xz/lzma"
)

func compress(ext string) func(io.Writer) (io.WriteCloser, error) {
	switch ext {
	case ".gz":
		return func(r io.Writer) (io.WriteCloser, error) {
			return gzip.NewWriter(r), nil
		}
	case ".xz":
		return func(r io.Writer) (io.WriteCloser, error) {
			return xz.NewWriter(r)
		}
	case ".lzma":
		return func(r io.Writer) (io.WriteCloser, error) {
			return lzma.NewWriter(r)
		}
	case ".zst":
		return func(r io.Writer) (io.WriteCloser, error) {
			return zstd.NewWriter(r)
		}
	default:
		return func(r io.Writer) (io.WriteCloser, error) {
			return writeCloser{r}, nil
		}
	}
}

func decompress(ext string) func(io.Reader) (io.ReadCloser, error) {
	switch ext {
	case ".gz":
		return func(r io.Reader) (io.ReadCloser, error) {
			return gzip.NewReader(r)
		}
	case ".bz2":
		return func(r io.Reader) (io.ReadCloser, error) {
			return io.NopCloser(bzip2.NewReader(r)), nil
		}
	case ".xz":
		return func(r io.Reader) (io.ReadCloser, error) {
			rd, err := xz.NewReader(r)
			return io.NopCloser(rd), err
		}
	case ".lzma":
		return func(r io.Reader) (io.ReadCloser, error) {
			rd, err := lzma.NewReader(r)
			return io.NopCloser(rd), err
		}
	case ".zst":
		return func(r io.Reader) (io.ReadCloser, error) {
			rd, err := zstd.NewReader(r)
			return readCloser{rd}, err
		}
	default:
		return func(r io.Reader) (io.ReadCloser, error) {
			return io.NopCloser(r), nil
		}
	}
}

type writeCloser struct {
	io.Writer
}

func (c writeCloser) Close() error {
	return nil
}

type rc interface {
	io.Reader
	Close()
}

type readCloser struct {
	rc
}

func (c readCloser) Close() error {
	c.rc.Close()
	return nil
}
