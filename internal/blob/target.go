package blob

import (
	"context"
	"io"
	"mime"
	"net/url"
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

func NewTarget(url, prefix, baseURL string) (target.Target, error) {
	b, err := blob.OpenBucket(context.Background(), url)
	if err != nil {
		return nil, err
	}

	t := &blobTarget{
		bucket:  b,
		prefix:  strings.Trim(prefix, "/"),
		baseURL: strings.TrimRight(baseURL, "/"),
	}

	return t, nil
}

func (t *blobTarget) NewWriter(ctx context.Context, filename string) (io.WriteCloser, error) {
	opt := &blob.WriterOptions{ContentType: mime.TypeByExtension(path.Ext(filename))}
	return t.bucket.NewWriter(ctx, path.Join(t.prefix, filename), opt)
}

func (t *blobTarget) NewReader(ctx context.Context, filename string) (io.ReadCloser, error) {
	return t.bucket.NewReader(ctx, path.Join(t.prefix, filename), nil)
}

func (t *blobTarget) Sub(dir string) target.Target {
	return &blobTarget{bucket: t.bucket, prefix: path.Join(t.prefix, strings.Trim(dir, "/")), baseURL: t.baseURL}
}

func (t *blobTarget) URL(ctx context.Context, filename string) (string, error) {
	if t.baseURL != "" {
		return url.JoinPath(t.baseURL, t.prefix, filename)
	}

	url, err := t.bucket.SignedURL(ctx, path.Join(t.prefix, filename), &blob.SignedURLOptions{})
	if err != nil {
		return "", err
	}

	url, _, _ = strings.Cut(url, "?")
	return url, nil
}
