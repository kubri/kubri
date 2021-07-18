package source

import (
	"fmt"
	"log"
	"net/url"
	"os"
	"strings"
	"time"

	"github.com/abemedia/appcast/pkg/slices"
	"github.com/hashicorp/go-version"
	"golang.org/x/mod/semver"
)

type Config struct {
	Token string
	Repo  string
}

type Factory = func(Config) (*Source, error)

//nolint:gochecknoglobals
var sources = map[string]Factory{}

func Register(scheme string, factory Factory) {
	sources[scheme] = factory
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

type Provider interface {
	GetRelease(version string) (*Release, error)
	ListReleases() ([]*Release, error)
	DownloadAsset(version, name string) ([]byte, error)
	UploadAsset(version, name string, data []byte) error
}

type Source struct {
	Provider
}

type ListOptions struct {
	Constraint string
}

func (s *Source) ListReleases(opt *ListOptions) ([]*Release, error) {
	var constraint version.Constraints
	if opt != nil && opt.Constraint != "" {
		c, err := version.NewConstraint(opt.Constraint)
		if err != nil {
			return nil, err
		}
		constraint = c
	}

	releases, err := s.Provider.ListReleases()
	if err != nil {
		return nil, err
	}

	releases = slices.Filter(releases, func(r *Release) bool {
		v, err := version.NewSemver(r.Version)
		if err != nil {
			log.Println("Skipping invalid version:", r.Version)
			return false
		}

		if !constraint.Check(v) {
			return false
		}

		processRelease(r)

		return true
	})

	return releases, nil
}

func (s *Source) GetRelease(version string) (*Release, error) {
	r, err := s.Provider.GetRelease(version)
	if err != nil {
		return nil, err
	}

	processRelease(r)

	return r, nil
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

func processRelease(r *Release) {
	if r.Name == "" {
		r.Name = r.Version
	}
	if !r.Prerelease {
		r.Prerelease = semver.Prerelease(r.Version) != ""
	}
}
