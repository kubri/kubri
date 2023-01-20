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
	bucket  *blob.Bucket
	prefix  string
	baseURL string
}

func New(url, prefix, baseURL string) (target.Target, error) {
	b, err := blob.OpenBucket(context.Background(), url)
	if err != nil {
		return nil, err
	}

	if prefix != "" {
		prefix = path.Clean(prefix)
	}

	return &blobTarget{bucket: b, prefix: prefix, baseURL: baseURL}, nil
}

func (s *blobTarget) NewWriter(ctx context.Context, filename string) (io.WriteCloser, error) {
	opt := &blob.WriterOptions{ContentType: mime.TypeByExtension(path.Ext(filename))}
	return s.bucket.NewWriter(ctx, path.Join(s.prefix, filename), opt)
}

func (s *blobTarget) NewReader(ctx context.Context, filename string) (io.ReadCloser, error) {
	return s.bucket.NewReader(ctx, path.Join(s.prefix, filename), nil)
}

func (s *blobTarget) Sub(dir string) target.Target {
	return &blobTarget{bucket: s.bucket, prefix: path.Join(s.prefix, dir)}
}

func (s *blobTarget) URL(ctx context.Context, filename string) (string, error) {
	if s.baseURL != "" {
		return s.baseURL + "/" + path.Join(s.prefix, filename), nil
	}

	url, err := s.bucket.SignedURL(ctx, path.Join(s.prefix, filename), &blob.SignedURLOptions{})
	if err != nil {
		return "", err
	}

	url, _, _ = strings.Cut(url, "?")
	return url, nil
}
