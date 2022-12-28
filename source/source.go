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

type ByVersion []*Release

func (vs ByVersion) Len() int      { return len(vs) }
func (vs ByVersion) Swap(i, j int) { vs[i], vs[j] = vs[j], vs[i] }
func (vs ByVersion) Less(i, j int) bool {
	cmp := semver.Compare(vs[i].Version, vs[j].Version)
	if cmp != 0 {
		return cmp > 0
	}
	return vs[i].Version > vs[j].Version
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

type ListOptions struct {
	Constraint string
}

func (s *Source) ListReleases(ctx context.Context, opt *ListOptions) ([]*Release, error) {
	if s == nil || s.s == nil {
		return nil, ErrMissingSource
	}

	var constraint version.Constraints
	if opt != nil && opt.Constraint != "" {
		c, err := version.NewConstraint(opt.Constraint)
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
			return false
		}

		processRelease(r)

		return true
	})

	sort.Sort(ByVersion(releases))

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

func (s *Source) UnmarshalText(b []byte) error {
	provider, repo, ok := strings.Cut(string(b), "://")
	if !ok {
		return fmt.Errorf("invalid source URL: %s", b)
	}

	factory, ok := sources[provider]
	if !ok {
		return fmt.Errorf("unsupported source: %s", provider)
	}

	source, err := factory(Config{
		Repo:  repo,
		Token: os.Getenv(strings.ToUpper(provider) + "_TOKEN"),
	})
	if err != nil {
		return err
	}

	s.s = source.s

	return nil
}

func processRelease(r *Release) {
	if r.Name == "" {
		r.Name = r.Version
	}
	if !r.Prerelease {
		r.Prerelease = semver.Prerelease(r.Version) != ""
	}
}
