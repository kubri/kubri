package gitlab

import (
	"bytes"
	"net/http"
	"strconv"

	"github.com/abemedia/appcast/source"
	"github.com/hashicorp/go-retryablehttp"
	"github.com/xanzy/go-gitlab"
)

type gitlabSource struct {
	client *gitlab.Client
	repo   string
}

func New(c source.Config) (*source.Source, error) {
	s := new(gitlabSource)

	git, err := gitlab.NewClient(c.Token)
	if err != nil {
		return nil, err
	}

	s.client = git
	s.repo = c.Repo

	return &source.Source{Provider: s}, nil
}

func (s *gitlabSource) Releases() ([]*source.Release, error) {
	releases, _, err := s.client.Releases.ListReleases(s.repo, nil)
	if err != nil {
		return nil, err
	}

	result := make([]*source.Release, 0, len(releases))
	for _, release := range releases {
		r := &source.Release{
			Name:        release.Name,
			Description: release.Description,
			Version:     release.TagName,
			Date:        *release.CreatedAt,
			Assets:      make([]*source.Asset, 0, len(release.Assets.Links)),
		}

		for _, l := range release.Assets.Links {
			size, err := s.getSize(l.URL)
			if err != nil {
				return nil, err
			}

			r.Assets = append(r.Assets, &source.Asset{
				Name: l.Name,
				URL:  l.URL,
				Size: size,
			})
		}

		result = append(result, r)
	}

	return result, nil
}

func (s *gitlabSource) UploadAsset(version, name string, data []byte) error {
	file, _, err := s.client.Projects.UploadFile(s.repo, bytes.NewReader(data), name)
	if err != nil {
		return err
	}

	_, _, err = s.client.ReleaseLinks.CreateReleaseLink(s.repo, version, &gitlab.CreateReleaseLinkOptions{
		Name: &name,
		URL:  &file.URL,
	})

	return err
}

func (s *gitlabSource) DownloadAsset(version, name string) ([]byte, error) {
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
			_, err = s.client.Do(req, &buf)
			if err != nil {
				return nil, err
			}

			return buf.Bytes(), nil
		}
	}

	return nil, source.ErrAssetNotFound
}

func (s *gitlabSource) getSize(url string) (int, error) {
	req, err := retryablehttp.NewRequest(http.MethodHead, url, nil)
	if err != nil {
		return 0, err
	}

	r, err := s.client.Do(req, nil)
	if err != nil {
		return 0, err
	}

	return strconv.Atoi(r.Header.Get("Content-Length"))
}

//nolint:gochecknoinits
func init() { source.Register("gitlab", New) }
