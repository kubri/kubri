package github

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"os"
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

	return &source.Source{Provider: s}, nil
}

func (s *githubSource) ListReleases() ([]*source.Release, error) {
	releases, _, err := s.client.Repositories.ListReleases(context.Background(), s.org, s.repo, nil)
	if err != nil {
		return nil, err
	}

	r := make([]*source.Release, 0, len(releases))
	for _, release := range releases {
		r = append(r, parseRelease(release))
	}

	return r, nil
}

func (s *githubSource) GetRelease(version string) (*source.Release, error) {
	r, _, err := s.client.Repositories.GetReleaseByTag(context.Background(), s.org, s.repo, version)
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

func (s *githubSource) UploadAsset(version, name string, data []byte) error {
	f, err := os.CreateTemp("", "")
	if err != nil {
		return err
	}

	_, err = f.Write(data)
	if err != nil {
		return err
	}

	ctx := context.Background()
	release, _, err := s.client.Repositories.GetReleaseByTag(ctx, s.org, s.repo, version)
	if err != nil {
		return err
	}

	opt := &github.UploadOptions{Name: name}
	_, _, err = s.client.Repositories.UploadReleaseAsset(ctx, s.org, s.repo, release.GetID(), opt, f)
	if err != nil {
		return err
	}

	return nil
}

func (s *githubSource) DownloadAsset(version, name string) ([]byte, error) {
	ctx := context.Background()
	release, _, err := s.client.Repositories.GetReleaseByTag(ctx, s.org, s.repo, version)
	if err != nil {
		return nil, err
	}

	for _, asset := range release.Assets {
		if asset.GetName() == name {
			rc, _, err := s.client.Repositories.DownloadReleaseAsset(ctx, s.org, s.repo, release.GetID())
			if err != nil {
				return nil, err
			}
			defer rc.Close()

			return io.ReadAll(rc)
		}
	}

	return nil, source.ErrAssetNotFound
}

//nolint:gochecknoinits
func init() { source.Register("github", New) }
