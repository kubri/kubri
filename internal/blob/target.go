package blob

import (
	"context"
	"io"
	"io/fs"
	"mime"
	"net/url"
	"path"
	"strings"

	"gocloud.dev/blob"
	"gocloud.dev/gcerrors"

	"github.com/kubri/kubri/target"
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
	f, err := t.bucket.NewWriter(ctx, path.Join(t.prefix, filename), opt)
	return f, mapError("write", filename, err)
}

func (t *blobTarget) NewReader(ctx context.Context, filename string) (io.ReadCloser, error) {
	f, err := t.bucket.NewReader(ctx, path.Join(t.prefix, filename), nil)
	return f, mapError("read", filename, err)
}

func (t *blobTarget) Remove(ctx context.Context, filename string) error {
	return mapError("remove", filename, t.bucket.Delete(ctx, path.Join(t.prefix, filename)))
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
		return "", mapError("read", filename, err)
	}

	url, _, _ = strings.Cut(url, "?")
	return url, nil
}

func mapError(op, name string, err error) error {
	switch gcerrors.Code(err) {
	case gcerrors.OK:
		return nil
	case gcerrors.NotFound:
		return &fs.PathError{Op: op, Path: name, Err: fs.ErrNotExist}
	default:
		return err
	}
}
