package blob

import (
	"context"
	"io"
	"mime"
	"path"
	"strings"

	"github.com/abemedia/appcast/target"
	"gocloud.dev/blob"
)

type blobTarget struct {
	bucket *blob.Bucket
	prefix string
}

func New(url, prefix string) (target.Target, error) {
	b, err := blob.OpenBucket(context.Background(), url)
	if err != nil {
		return nil, err
	}

	if prefix != "" {
		prefix = strings.Trim(prefix, "/") + "/"
	}

	return &blobTarget{bucket: b, prefix: prefix}, nil
}

func (s *blobTarget) NewWriter(ctx context.Context, filename string) (io.WriteCloser, error) {
	opt := &blob.WriterOptions{ContentType: mime.TypeByExtension(path.Ext(filename))}
	return s.bucket.NewWriter(ctx, path.Join(s.prefix, filename), opt)
}

func (s *blobTarget) NewReader(ctx context.Context, filename string) (io.ReadCloser, error) {
	return s.bucket.NewReader(ctx, path.Join(s.prefix, filename), nil)
}

func (s *blobTarget) Sub(dir string) target.Target {
	sub := *s
	sub.prefix += strings.Trim(dir, "/") + "/"
	return &sub
}
