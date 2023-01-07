package pipe

import (
	"context"
	"errors"
	"os"
	"path/filepath"

	"github.com/abemedia/appcast/integrations/apt"
	"github.com/abemedia/appcast/integrations/sparkle"
	"github.com/abemedia/appcast/source"
	"github.com/abemedia/appcast/target"
	"gopkg.in/yaml.v3"
)

type config struct {
	Title       string `yaml:"title"`
	Description string `yaml:"description"`
	Source      struct {
		Repo       string `yaml:"repo"`
		Prerelease bool   `yaml:"prerelease"`
		Version    string `yaml:"version"`
	} `yaml:"source"`
	Target struct {
		Repo string `yaml:"repo"`
		Flat bool   `yaml:"flat"`
	} `yaml:"target"`
	Apt     aptConfig     `yaml:"apt"`
	Sparkle sparkleConfig `yaml:"sparkle"`

	source *source.Source
	target target.Target
}

type Pipe struct {
	Apt     *apt.Config
	Sparkle *sparkle.Config
}

func (p *Pipe) Run(ctx context.Context) error {
	var err error
	if p.Apt != nil {
		if err = apt.Build(ctx, p.Apt); err != nil {
			return err
		}
	}
	if p.Sparkle != nil {
		if err = sparkle.Build(ctx, p.Sparkle); err != nil {
			return err
		}
	}
	return nil
}

func Load(path string) (*Pipe, error) {
	b, err := open(path)
	if err != nil {
		return nil, err
	}

	c := &config{}
	if err = yaml.Unmarshal(b, c); err != nil {
		return nil, err
	}
	if c.source, err = source.Open(c.Source.Repo); err != nil {
		return nil, err
	}
	if c.target, err = target.Open(c.Target.Repo); err != nil {
		return nil, err
	}

	var p Pipe
	if !c.Apt.Disabled {
		p.Apt = getApt(c)
	}
	if !c.Sparkle.Disabled {
		if p.Sparkle, err = getSparkle(c); err != nil {
			return nil, err
		}
	}

	return &p, nil
}

func open(path string) ([]byte, error) {
	if path != "" {
		return os.ReadFile(path)
	}

	paths := []string{
		".appcast.yml",
		".appcast.yaml",
		"appcast.yml",
		"appcast.yaml",
		filepath.Join(".github", "appcast.yml"),
		filepath.Join(".github", "appcast.yaml"),
	}

	for _, path := range paths {
		b, err := os.ReadFile(path)
		if err != nil {
			if os.IsNotExist(err) {
				continue
			}
			return nil, err
		}
		return b, nil
	}

	return nil, errors.New("no config file found")
}

func fallback[T comparable](a, b T) T {
	var zero T
	if a != zero {
		return a
	}
	return b
}
