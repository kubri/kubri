package blob

import (
	"context"
	"io"
	"mime"
	"path"
	"strings"

	"github.com/abemedia/appcast/source"
	"gocloud.dev/blob"
)

type blobSource struct {
	bucket  *blob.Bucket
	prefix  string
	baseURL string
}

func New(url, prefix, baseURL string) (*source.Source, error) {
	b, err := blob.OpenBucket(context.Background(), url)
	if err != nil {
		return nil, err
	}

	if prefix != "" {
		prefix = strings.Trim(prefix, "/") + "/"
	}

	s := &blobSource{
		bucket:  b,
		prefix:  prefix,
		baseURL: baseURL,
	}

	return source.New(s), nil
}

func (s *blobSource) ListReleases(ctx context.Context) ([]*source.Release, error) {
	releases := []*source.Release{}
	iter := s.bucket.List(&blob.ListOptions{Prefix: s.prefix, Delimiter: "/"})

	for {
		object, err := iter.Next(ctx)
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, err
		}

		if !object.IsDir {
			continue
		}

		r, err := s.GetRelease(ctx, path.Base(object.Key))
		if err != nil && err != source.ErrNoReleaseFound {
			return nil, err
		}

		if r != nil {
			releases = append(releases, r)
		}
	}

	return releases, nil
}

func (s *blobSource) GetRelease(ctx context.Context, version string) (*source.Release, error) {
	r := source.Release{Version: version}

	iter := s.bucket.List(&blob.ListOptions{Prefix: path.Join(s.prefix, version) + "/", Delimiter: "/"})
	for {
		object, err := iter.Next(ctx)
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, err
		}

		if object.IsDir {
			continue
		}

		if r.Date.IsZero() || r.Date.After(object.ModTime) {
			r.Date = object.ModTime
		}

		url, err := s.bucket.SignedURL(ctx, object.Key, &blob.SignedURLOptions{})
		if err != nil {
			url = s.baseURL + "/" + object.Key
		} else {
			url, _, _ = strings.Cut(url, "?")
		}

		attr, err := s.bucket.Attributes(ctx, object.Key)
		if err != nil {
			return nil, err
		}

		r.Assets = append(r.Assets, &source.Asset{
			Name: path.Base(object.Key),
			URL:  url,
			Size: int(attr.Size),
		})
	}

	if len(r.Assets) == 0 {
		return nil, source.ErrNoReleaseFound
	}

	return &r, nil
}

func (s *blobSource) UploadAsset(ctx context.Context, version, name string, data []byte) error {
	opt := &blob.WriterOptions{ContentType: mime.TypeByExtension(path.Ext(name))}
	return s.bucket.WriteAll(ctx, path.Join(s.prefix, version, name), data, opt)
}

func (s *blobSource) DownloadAsset(ctx context.Context, version, name string) ([]byte, error) {
	return s.bucket.ReadAll(ctx, path.Join(s.prefix, version, name))
}
