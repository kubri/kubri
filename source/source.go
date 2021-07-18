package source

import (
	"fmt"
	"net/url"
	"os"
	"strings"
	"time"

	"golang.org/x/mod/semver"
)

type Factory = func(Config) (*Source, error)

//nolint:gochecknoglobals
var sources = map[string]Factory{}

func Register(scheme string, factory Factory) {
	sources[scheme] = factory
}

type Config struct {
	Token string
	Repo  string
}

type Provider interface {
	Releases() ([]*Release, error)
	DownloadAsset(version, name string) ([]byte, error)
	UploadAsset(version, name string, data []byte) error
}

type Source struct {
	Provider
}

func (s *Source) Releases() ([]*Release, error) {
	releases, err := s.Provider.Releases()
	if err != nil {
		return nil, err
	}
	for _, r := range releases {
		if r.Name == "" {
			r.Name = r.Version
		}
		if !r.Prerelease {
			r.Prerelease = semver.Prerelease(r.Version) != ""
		}
	}
	return releases, nil
}

func (s *Source) UnmarshalText(b []byte) error {
	u, err := url.Parse(string(b))
	if err != nil {
		return err
	}

	for scheme, factory := range sources {
		if u.Scheme == scheme {
			source, err := factory(Config{
				Repo:  u.Host + u.Path,
				Token: os.Getenv(strings.ToUpper(scheme) + "_TOKEN"),
			})
			if err != nil {
				return err
			}
			*s = *source
			return nil
		}
	}

	return fmt.Errorf("unsupported source: %s", u.Scheme)
}

type Release struct {
	Name        string
	Description string
	Date        time.Time
	Version     string
	Prerelease  bool
	Assets      []*Asset
}

type Asset struct {
	Name string
	URL  string
	Size int
}
