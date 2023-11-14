package blob

import (
	"context"
	"errors"
	"io"
	"math"
	"mime"
	"net/url"
	"path"
	"strings"
	"time"

	"github.com/abemedia/appcast/source"
	"github.com/parkr/changelog"
	"gocloud.dev/blob"
)

type blobSource struct {
	bucket  *blob.Bucket
	prefix  string
	baseURL string
}

func NewSource(url, prefix, baseURL string) (*source.Source, error) {
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

		if strings.EqualFold(path.Base(object.Key), "CHANGELOG.md") {
			rd, err := s.bucket.NewReader(ctx, object.Key, nil)
			if err != nil {
				return nil, err
			}
			date, desc, err := parseChangelog(rd, version)
			if err != nil {
				return nil, err
			}
			if r.Date.IsZero() {
				r.Date = date
			}
			r.Description = desc
			continue
		}

		var u string
		if s.baseURL != "" {
			u, err = url.JoinPath(s.baseURL, object.Key)
			if err != nil {
				return nil, err
			}
		} else {
			u, err = s.bucket.SignedURL(ctx, object.Key, &blob.SignedURLOptions{Expiry: math.MaxInt64})
			if err != nil {
				return nil, err
			}
		}

		attr, err := s.bucket.Attributes(ctx, object.Key)
		if err != nil {
			return nil, err
		}

		r.Assets = append(r.Assets, &source.Asset{
			Name: path.Base(object.Key),
			URL:  u,
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

func parseChangelog(rd io.Reader, version string) (time.Time, string, error) {
	var date time.Time

	c, err := changelog.NewChangelogFromReader(rd)
	if err != nil {
		return time.Time{}, "", err
	}
	v := c.GetVersion(strings.TrimPrefix(version, "v"))
	if v == nil {
		return time.Time{}, "", errors.New("changelog doesn't contain " + version)
	}

	if v.Date != "" {
		d, err := time.Parse(time.DateOnly, v.Date)
		if err != nil {
			return time.Time{}, "", err
		}
		date = d
	}

	var w strings.Builder
	w.Grow(512)

	if len(v.History) > 0 {
		w.WriteString(v.History[0].Summary)
		for _, history := range v.History[1:] {
			w.WriteString("* ")
			w.WriteString(history.Summary)
		}
	}

	if len(v.Subsections) > 0 {
		for i, subsection := range v.Subsections {
			if i != 0 {
				w.WriteString("\n\n")
			}
			w.WriteString("### ")
			w.WriteString(subsection.Name)
			w.WriteString("\n\n")
			w.WriteString(subsection.History[0].Summary)
			for _, history := range subsection.History[1:] {
				w.WriteString("* ")
				w.WriteString(history.Summary)
			}
		}
	}

	return date, w.String(), nil
}
