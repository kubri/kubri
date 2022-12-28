package github

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"mime"
	"net/http"
	"path"
	"strings"

	"github.com/abemedia/appcast/source"
	"github.com/google/go-github/github"
	"golang.org/x/oauth2"
)

type githubSource struct {
	client *github.Client
	org    string
	repo   string
}

func New(c source.Config) (*source.Source, error) {
	org, repo, ok := strings.Cut(c.Repo, "/")
	if !ok {
		return nil, fmt.Errorf("invalid repo: %s", c.Repo)
	}

	var client *http.Client
	if c.Token != "" {
		ts := oauth2.StaticTokenSource(&oauth2.Token{AccessToken: c.Token})
		client = oauth2.NewClient(context.Background(), ts)
	}

	s := &githubSource{
		client: github.NewClient(client),
		org:    org,
		repo:   repo,
	}

	return source.New(s), nil
}

func (s *githubSource) ListReleases(ctx context.Context) ([]*source.Release, error) {
	releases, _, err := s.client.Repositories.ListReleases(ctx, s.org, s.repo, nil)
	if err != nil {
		return nil, err
	}

	r := make([]*source.Release, 0, len(releases))
	for _, release := range releases {
		r = append(r, parseRelease(release))
	}

	return r, nil
}

func (s *githubSource) GetRelease(ctx context.Context, version string) (*source.Release, error) {
	r, _, err := s.client.Repositories.GetReleaseByTag(ctx, s.org, s.repo, version)
	if err != nil {
		return nil, err
	}

	return parseRelease(r), nil
}

func parseRelease(release *github.RepositoryRelease) *source.Release {
	r := &source.Release{
		Name:        release.GetName(),
		Description: release.GetBody(),
		Version:     release.GetTagName(),
		Date:        release.PublishedAt.Time,
		Prerelease:  release.GetPrerelease(),
		Assets:      make([]*source.Asset, 0, len(release.Assets)),
	}

	for _, asset := range release.Assets {
		r.Assets = append(r.Assets, &source.Asset{
			Name: asset.GetName(),
			URL:  asset.GetBrowserDownloadURL(),
			Size: asset.GetSize(),
		})
	}

	return r
}

func (s *githubSource) UploadAsset(ctx context.Context, version, name string, data []byte) error {
	release, _, err := s.client.Repositories.GetReleaseByTag(ctx, s.org, s.repo, version)
	if err != nil {
		return err
	}

	u := fmt.Sprintf("repos/%s/%s/releases/%d/assets?name=%s", s.org, s.repo, release.GetID(), name)
	mediaType := mime.TypeByExtension(path.Ext(name))
	req, err := s.client.NewUploadRequest(u, bytes.NewReader(data), int64(len(data)), mediaType)
	if err != nil {
		return err
	}

	_, err = s.client.Do(ctx, req, nil)
	return err
}

func (s *githubSource) DownloadAsset(ctx context.Context, version, name string) ([]byte, error) {
	release, _, err := s.client.Repositories.GetReleaseByTag(ctx, s.org, s.repo, version)
	if err != nil {
		return nil, err
	}

	for _, asset := range release.Assets {
		if asset.GetName() == name {
			r, u, err := s.client.Repositories.DownloadReleaseAsset(ctx, s.org, s.repo, asset.GetID())
			if err != nil {
				return nil, err
			}

			if r != nil {
				defer r.Close()
				return io.ReadAll(r)
			}

			req, err := http.NewRequestWithContext(ctx, http.MethodGet, u, nil)
			if err != nil {
				return nil, err
			}

			res, err := http.DefaultClient.Do(req)
			if err != nil {
				return nil, err
			}

			defer res.Body.Close()
			return io.ReadAll(res.Body)
		}
	}

	return nil, source.ErrAssetNotFound
}

//nolint:gochecknoinits
func init() { source.Register("github", New) }
