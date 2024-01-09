package pipe

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/abemedia/appcast/integrations/apk"
	"github.com/abemedia/appcast/integrations/appinstaller"
	"github.com/abemedia/appcast/integrations/apt"
	"github.com/abemedia/appcast/integrations/sparkle"
	"github.com/abemedia/appcast/integrations/yum"
	"github.com/abemedia/appcast/source"
	"github.com/abemedia/appcast/target"
	"golang.org/x/sync/errgroup"
	"gopkg.in/yaml.v3"
)

type config struct {
	Title          string              `yaml:"title,omitempty"`
	Description    string              `yaml:"description,omitempty"`
	Version        string              `yaml:"version,omitempty"`
	Prerelease     bool                `yaml:"prerelease,omitempty"`
	UploadPackages bool                `yaml:"upload-packages,omitempty"`
	Source         *sourceConfig       `yaml:"source"`
	Target         *targetConfig       `yaml:"target"`
	Apk            *apkConfig          `yaml:"apk,omitempty"`
	Apt            *aptConfig          `yaml:"apt,omitempty"`
	Yum            *yumConfig          `yaml:"yum,omitempty"`
	Sparkle        *sparkleConfig      `yaml:"sparkle,omitempty"`
	Appinstaller   *appinstallerConfig `yaml:"appinstaller,omitempty"`

	source *source.Source
	target target.Target
}

type Pipe struct {
	Apk          *apk.Config
	Appinstaller *appinstaller.Config
	Apt          *apt.Config
	Yum          *yum.Config
	Sparkle      *sparkle.Config
}

//nolint:funlen
func (p *Pipe) Run(ctx context.Context) error {
	var n int
	start := time.Now()
	g, ctx := errgroup.WithContext(ctx)

	if p.Apk != nil {
		n++
		g.Go(func() error {
			log.Print("Publishing APK packages...")
			if err := apk.Build(ctx, p.Apk); err != nil {
				return fmt.Errorf("failed to publish APK packages: %w", err)
			}
			log.Print("Completed publishing APK packages.")
			return nil
		})
	}
	if p.Appinstaller != nil {
		n++
		g.Go(func() error {
			log.Print("Publishing App Installer packages...")
			if err := appinstaller.Build(ctx, p.Appinstaller); err != nil {
				return fmt.Errorf("failed to publish App Installer packages: %w", err)
			}
			log.Print("Completed publishing App Installer packages.")
			return nil
		})
	}
	if p.Apt != nil {
		n++
		g.Go(func() error {
			log.Print("Publishing APT packages...")
			if err := apt.Build(ctx, p.Apt); err != nil {
				return fmt.Errorf("failed to publish APT packages: %w", err)
			}
			log.Print("Completed publishing APT packages.")
			return nil
		})
	}
	if p.Yum != nil {
		n++
		g.Go(func() error {
			log.Print("Publishing YUM packages...")
			if err := yum.Build(ctx, p.Yum); err != nil {
				return fmt.Errorf("failed to publish YUM packages: %w", err)
			}
			log.Print("Completed publishing YUM packages.")
			return nil
		})
	}
	if p.Sparkle != nil {
		n++
		g.Go(func() error {
			log.Print("Publishing Sparkle packages...")
			if err := sparkle.Build(ctx, p.Sparkle); err != nil {
				return fmt.Errorf("failed to publish Sparkle packages: %w", err)
			}
			log.Print("Completed publishing Sparkle packages.")
			return nil
		})
	}

	if n == 0 {
		return errors.New("no integrations configured")
	}

	if err := g.Wait(); err != nil {
		return err
	}

	log.Printf("Completed in %s", time.Since(start).Truncate(time.Millisecond))

	return nil
}

func Load(path string) (*Pipe, error) {
	b, err := open(path)
	if err != nil {
		return nil, err
	}

	dec := yaml.NewDecoder(bytes.NewReader(b))
	dec.KnownFields(true)

	c := &config{}
	if err = dec.Decode(c); err != nil {
		return nil, err
	}
	if c.source, err = getSource(c.Source); err != nil {
		return nil, err
	}
	if c.target, err = getTarget(c.Target); err != nil {
		return nil, err
	}

	var p Pipe
	if c.Apk != nil && !c.Apk.Disabled {
		if p.Apk, err = getApk(c); err != nil {
			return nil, err
		}
	}
	if c.Appinstaller != nil && !c.Appinstaller.Disabled {
		p.Appinstaller = getAppinstaller(c)
	}
	if c.Apt != nil && !c.Apt.Disabled {
		if p.Apt, err = getApt(c); err != nil {
			return nil, err
		}
	}
	if c.Yum != nil && !c.Yum.Disabled {
		if p.Yum, err = getYum(c); err != nil {
			return nil, err
		}
	}
	if c.Sparkle != nil && !c.Sparkle.Disabled {
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
