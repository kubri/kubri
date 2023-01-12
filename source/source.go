package source

import (
	"context"
	"fmt"
	"log"
	"os"
	"sort"
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

var providers = map[string]Factory{} //nolint:gochecknoglobals

func Register(scheme string, factory Factory) {
	providers[scheme] = factory
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

type Driver interface {
	GetRelease(ctx context.Context, version string) (*Release, error)
	ListReleases(ctx context.Context) ([]*Release, error)
	DownloadAsset(ctx context.Context, version, name string) ([]byte, error)
	UploadAsset(ctx context.Context, version, name string, data []byte) error
}

type Source struct {
	s Driver
}

func New(driver Driver) *Source {
	return &Source{s: driver}
}

func Open(url string) (*Source, error) {
	provider, repo, ok := strings.Cut(url, "://")
	if !ok {
		return nil, fmt.Errorf("invalid source URL: %s", url)
	}

	factory, ok := providers[provider]
	if !ok {
		return nil, fmt.Errorf("unsupported source: %s", provider)
	}

	return factory(Config{
		Repo:  repo,
		Token: os.Getenv(strings.ToUpper(provider) + "_TOKEN"),
	})
}

type ListOptions struct {
	// Version constraint e.g. 'v1.2.4', 'v1', '>= v1.1.0, < v2.1'
	Version string

	// Include pre-releases
	Prerelease bool
}

func (s *Source) ListReleases(ctx context.Context, opt *ListOptions) ([]*Release, error) {
	if s == nil || s.s == nil {
		return nil, ErrMissingSource
	}

	var constraint version.Constraints
	if opt != nil && opt.Version != "" && opt.Version != "latest" {
		c, err := version.NewConstraint(opt.Version)
		if err != nil {
			return nil, err
		}
		constraint = c
	}

	releases, err := s.s.ListReleases(ctx)
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
			log.Println("Skipping version:", r.Version)
			return false
		}

		if (opt == nil || !opt.Prerelease) && v.Prerelease() != "" {
			log.Println("Skipping prerelease:", r.Version)
			return false
		}

		processRelease(r)

		return true
	})

	sort.Sort(ByVersion(releases))

	if opt != nil && opt.Version == "latest" {
		return releases[:1], nil
	}

	return releases, nil
}

func (s *Source) GetRelease(ctx context.Context, version string) (*Release, error) {
	if s == nil || s.s == nil {
		return nil, ErrMissingSource
	}

	r, err := s.s.GetRelease(ctx, version)
	if err != nil {
		return nil, err
	}

	processRelease(r)

	return r, nil
}

func (s *Source) DownloadAsset(ctx context.Context, version, name string) ([]byte, error) {
	if s == nil || s.s == nil {
		return nil, ErrMissingSource
	}

	return s.s.DownloadAsset(ctx, version, name)
}

func (s *Source) UploadAsset(ctx context.Context, version, name string, data []byte) error {
	if s == nil || s.s == nil {
		return ErrMissingSource
	}

	return s.s.UploadAsset(ctx, version, name, data)
}

func processRelease(r *Release) {
	if r.Name == "" {
		r.Name = r.Version
	}
	if !r.Prerelease {
		r.Prerelease = semver.Prerelease(r.Version) != ""
	}
}
