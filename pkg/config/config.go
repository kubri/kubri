package config

import (
	"bytes"
	"errors"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"

	"github.com/kubri/kubri/integrations/apk"
	"github.com/kubri/kubri/integrations/appinstaller"
	"github.com/kubri/kubri/integrations/apt"
	"github.com/kubri/kubri/integrations/sparkle"
	"github.com/kubri/kubri/integrations/yum"
	"github.com/kubri/kubri/source"
	"github.com/kubri/kubri/target"
)

type Config struct {
	Apk          *apk.Config
	Appinstaller *appinstaller.Config
	Apt          *apt.Config
	Yum          *yum.Config
	Sparkle      *sparkle.Config
}

func Load(path string) (*Config, error) {
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

	if err = Validate(c); err != nil {
		return nil, err
	}

	if c.source, err = getSource(c.Source); err != nil {
		return nil, err
	}
	if c.target, err = getTarget(c.Target); err != nil {
		return nil, err
	}

	var p Config
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
		".kubri.yml",
		".kubri.yaml",
		"kubri.yml",
		"kubri.yaml",
		filepath.Join(".github", "kubri.yml"),
		filepath.Join(".github", "kubri.yaml"),
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

type config struct {
	Title          string              `yaml:"title,omitempty"`
	Description    string              `yaml:"description,omitempty"`
	Version        string              `yaml:"version,omitempty"         validate:"omitempty,version_constraint"`
	Prerelease     bool                `yaml:"prerelease,omitempty"`
	UploadPackages bool                `yaml:"upload-packages,omitempty"`
	Source         *sourceConfig       `yaml:"source"                    validate:"required"`
	Target         *targetConfig       `yaml:"target"                    validate:"required"`
	Apk            *apkConfig          `yaml:"apk,omitempty"`
	Apt            *aptConfig          `yaml:"apt,omitempty"`
	Yum            *yumConfig          `yaml:"yum,omitempty"`
	Sparkle        *sparkleConfig      `yaml:"sparkle,omitempty"`
	Appinstaller   *appinstallerConfig `yaml:"appinstaller,omitempty"`

	source *source.Source
	target target.Target
}

func fallback[T comparable](a, b T) T {
	var zero T
	if a != zero {
		return a
	}
	return b
}
