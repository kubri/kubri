package target

import (
	"context"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/abemedia/appcast/source"
)

type Factory = func(source.Config) (Target, error)

var providers = map[string]Factory{} //nolint:gochecknoglobals

func Register(scheme string, factory Factory) {
	providers[scheme] = factory
}

type Target interface {
	NewWriter(ctx context.Context, path string) (io.WriteCloser, error)
	NewReader(ctx context.Context, path string) (io.ReadCloser, error)
	Sub(dir string) Target
}

func Open(url string) (Target, error) {
	provider, repo, ok := strings.Cut(url, "://")
	if !ok {
		return nil, fmt.Errorf("invalid target URL: %s", url)
	}

	factory, ok := providers[provider]
	if !ok {
		return nil, fmt.Errorf("unsupported target: %s", provider)
	}

	return factory(source.Config{
		Repo:  repo,
		Token: os.Getenv(strings.ToUpper(provider) + "_TOKEN"),
	})
}
