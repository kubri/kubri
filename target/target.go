package target

import (
	"context"
	"io"
)

type Target interface {
	NewWriter(ctx context.Context, path string) (io.WriteCloser, error)
	NewReader(ctx context.Context, path string) (io.ReadCloser, error)
	Remove(ctx context.Context, path string) error
	Sub(dir string) Target
	URL(ctx context.Context, path string) (string, error)
}
