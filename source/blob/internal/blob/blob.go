package blob

import (
	"context"
	"io"
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
	ctx := context.Background()
	b, err := blob.OpenBucket(ctx, url)
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
	return &source.Source{Provider: s}, nil
}

func (s *blobSource) ListReleases() ([]*source.Release, error) {
	releases := []*source.Release{}
	ctx := context.Background()
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

		r, err := s.GetRelease(path.Base(object.Key))
		if err != nil && err != source.ErrReleaseNotFound {
			return nil, err
		}

		if r != nil {
			releases = append(releases, r)
		}
	}

	return releases, nil
}

func (s *blobSource) GetRelease(version string) (*source.Release, error) {
	ctx := context.Background()

	r := &source.Release{Version: version}

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
		return nil, source.ErrReleaseNotFound
	}

	return r, nil
}

func (s *blobSource) UploadAsset(version, name string, data []byte) error {
	return s.bucket.WriteAll(context.Background(), path.Join(s.prefix, version, name), data, nil)
}

func (s *blobSource) DownloadAsset(version, name string) ([]byte, error) {
	return s.bucket.ReadAll(context.Background(), path.Join(s.prefix, version, name))
}
