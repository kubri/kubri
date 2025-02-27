// Package gitlab provides a source implementation for GitLab releases.
package gitlab

import (
	"bytes"
	"context"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/hashicorp/go-retryablehttp"
	gitlab "gitlab.com/gitlab-org/api/client-go"

	"github.com/kubri/kubri/source"
)

// Config represents the configuration for a GitLab source.
type Config struct {
	Owner string
	Repo  string
	URL   string
}

// New returns a new GitLab source.
func New(c Config) (*source.Source, error) {
	var opt []gitlab.ClientOptionFunc
	if c.URL != "" {
		opt = append(opt, gitlab.WithBaseURL(c.URL))
	}

	client, err := gitlab.NewClient(os.Getenv("GITLAB_TOKEN"), opt...)
	if err != nil {
		return nil, err
	}

	s := &gitlabSource{
		client: client,
		repo:   c.Owner + "/" + c.Repo,
	}

	return source.New(s), nil
}

type gitlabSource struct {
	client *gitlab.Client
	repo   string
}

func (s *gitlabSource) ListReleases(ctx context.Context) ([]*source.Release, error) {
	releases, _, err := s.client.Releases.ListReleases(s.repo, nil, gitlab.WithContext(ctx))
	if err != nil {
		return nil, err
	}

	r := make([]*source.Release, 0, len(releases))
	for _, release := range releases {
		r = append(r, s.parseRelease(ctx, release))
	}

	return r, nil
}

func (s *gitlabSource) GetRelease(ctx context.Context, version string) (*source.Release, error) {
	r, _, err := s.client.Releases.GetRelease(s.repo, version, gitlab.WithContext(ctx))
	if err != nil {
		return nil, err
	}

	return s.parseRelease(ctx, r), nil
}

func (s *gitlabSource) parseRelease(ctx context.Context, release *gitlab.Release) *source.Release {
	r := &source.Release{
		Name:        release.Name,
		Description: release.Description,
		Version:     release.TagName,
		Date:        *release.CreatedAt,
		Assets:      make([]*source.Asset, 0, len(release.Assets.Links)),
	}

	for _, l := range release.Assets.Links {
		size, err := s.getSize(ctx, l.URL)
		if err != nil {
			log.Printf("failed to get size for %s: %s\n", l.Name, err)
		}

		r.Assets = append(r.Assets, &source.Asset{
			Name: l.Name,
			URL:  l.URL,
			Size: size,
		})
	}

	return r
}

func (s *gitlabSource) UploadAsset(ctx context.Context, version, name string, data []byte) error {
	file, _, err := s.client.ProjectMarkdownUploads.UploadProjectMarkdown(
		s.repo,
		bytes.NewReader(data),
		name,
		gitlab.WithContext(ctx),
	)
	if err != nil {
		return err
	}

	u := s.client.BaseURL()
	u.Path = file.FullPath
	url := u.String()

	opt := &gitlab.CreateReleaseLinkOptions{Name: &name, URL: &url}
	_, _, err = s.client.ReleaseLinks.CreateReleaseLink(s.repo, version, opt, gitlab.WithContext(ctx))

	return err
}

func (s *gitlabSource) DownloadAsset(ctx context.Context, version, name string) ([]byte, error) {
	links, _, err := s.client.ReleaseLinks.ListReleaseLinks(s.repo, version, nil)
	if err != nil {
		return nil, err
	}

	for _, link := range links {
		if link.Name == name {
			req, err := retryablehttp.NewRequest(http.MethodGet, link.URL, nil)
			if err != nil {
				return nil, err
			}

			var buf bytes.Buffer
			_, err = s.client.Do(req.WithContext(ctx), &buf)
			if err != nil {
				return nil, err
			}

			return buf.Bytes(), nil
		}
	}

	return nil, source.ErrAssetNotFound
}

func (s *gitlabSource) getSize(ctx context.Context, url string) (int, error) {
	req, err := retryablehttp.NewRequest(http.MethodHead, url, nil)
	if err != nil {
		return 0, err
	}

	r, err := s.client.Do(req.WithContext(ctx), nil)
	if err != nil {
		return 0, err
	}

	return strconv.Atoi(r.Header.Get("Content-Length"))
}
