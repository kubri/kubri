package config

import (
	"encoding/pem"
	"os"

	"github.com/abemedia/appcast/integrations/sparkle"
	"github.com/abemedia/appcast/pkg/crypto/dsa"
	"github.com/abemedia/appcast/pkg/crypto/ed25519"
	"github.com/abemedia/appcast/pkg/secret"
	"github.com/abemedia/appcast/source"
	"github.com/abemedia/appcast/target"
	"gopkg.in/yaml.v3"
)

type Config struct {
	Title       string  `yaml:"title"`
	Description string  `yaml:"description"`
	Source      Source  `yaml:"source"`
	Target      Target  `yaml:"target"`
	Sparkle     Sparkle `yaml:"sparkle"`
}

type Source struct {
	Repo       string `yaml:"repo"`
	Prerelease bool   `yaml:"prerelease"`
	Version    string `yaml:"version"`
}

type Target struct {
	Repo string `yaml:"repo"`
	Flat bool   `yaml:"flat"`
}

type Sparkle struct {
	Title       string         `yaml:"title"`
	Description string         `yaml:"description"`
	FileName    string         `yaml:"filename"`
	Params      []sparkle.Rule `yaml:"params"`
}

func Load(path string) (*Config, error) {
	b, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	c := &Config{}
	if err = yaml.Unmarshal(b, c); err != nil {
		return nil, err
	}
	return c, nil
}

func GetSparkle(c *Config) (*sparkle.Config, error) {
	src, err := source.Open(c.Source.Repo)
	if err != nil {
		return nil, err
	}

	tgt, err := target.Open(c.Target.Repo)
	if err != nil {
		return nil, err
	}
	if !c.Target.Flat {
		tgt = tgt.Sub("sparkle")
	}

	var dsaKey *dsa.PrivateKey
	if b, err := secret.Get("dsa_key"); err == nil {
		block, _ := pem.Decode(b)
		dsaKey, err = dsa.UnmarshalPrivateKey(block.Bytes)
		if err != nil {
			return nil, err
		}
	}

	var edKey ed25519.PrivateKey
	if b, err := secret.Get("ed25519_key"); err == nil {
		block, _ := pem.Decode(b)
		edKey, err = ed25519.UnmarshalPrivateKey(block.Bytes)
		if err != nil {
			return nil, err
		}
	}

	return &sparkle.Config{
		Title:       fallback(c.Sparkle.Title, c.Title),
		Description: fallback(c.Sparkle.Description, c.Description),
		FileName:    fallback(c.Sparkle.FileName, "sparkle.xml"),
		Source:      src,
		Target:      tgt,
		DSAKey:      dsaKey,
		Ed25519Key:  edKey,
		Version:     c.Source.Version,
		Prerelease:  c.Source.Prerelease,
		Settings:    c.Sparkle.Params,
	}, nil
}

func fallback[T comparable](a, b T) T {
	var zero T
	if a != zero {
		return a
	}
	return b
}
