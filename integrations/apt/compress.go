package apt

import (
	"compress/gzip"
	"io"

	"github.com/dsnet/compress/bzip2"
	"github.com/klauspost/compress/zstd"
	"github.com/pierrec/lz4"
	"github.com/ulikunitz/xz"
	"github.com/ulikunitz/xz/lzma"
)

type CompressionAlgo uint8

const (
	NoCompression CompressionAlgo = 1 << iota
	GZIP
	BZIP2
	XZ
	LZMA
	LZ4
	ZSTD
)

func compressionExtensions(algo CompressionAlgo) []string {
	if algo == 0 {
		return []string{"", ".xz", ".gz"} // default compression
	}

	a := []string{""} // Always start with blank string for uncompressed data.
	for i := GZIP; i <= ZSTD; i <<= 1 {
		if algo&i == 0 {
			continue
		}
		switch i {
		case NoCompression:
		case GZIP:
			a = append(a, ".gz")
		case BZIP2:
			a = append(a, ".bz2")
		case XZ:
			a = append(a, ".xz")
		case LZMA:
			a = append(a, ".lzma")
		case LZ4:
			a = append(a, ".lz4")
		case ZSTD:
			a = append(a, ".zst")
		default:
			panic("unknown compression algorithm")
		}
	}
	return a
}

func compress(ext string) func(io.Writer) (io.WriteCloser, error) {
	switch ext {
	case ".gz":
		return func(r io.Writer) (io.WriteCloser, error) {
			return gzip.NewWriter(r), nil
		}
	case ".bz2":
		return func(r io.Writer) (io.WriteCloser, error) {
			return bzip2.NewWriter(r, nil)
		}
	case ".xz":
		return func(r io.Writer) (io.WriteCloser, error) {
			return xz.NewWriter(r)
		}
	case ".lzma":
		return func(r io.Writer) (io.WriteCloser, error) {
			return lzma.NewWriter(r)
		}
	case ".lz4":
		return func(r io.Writer) (io.WriteCloser, error) {
			return lz4.NewWriter(r), nil
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
			return bzip2.NewReader(r, nil)
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
	case ".lz4":
		return func(r io.Reader) (io.ReadCloser, error) {
			return io.NopCloser(lz4.NewReader(r)), nil
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
