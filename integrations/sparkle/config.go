package sparkle

import (
	"github.com/abemedia/appcast/pkg/crypto/dsa"
	"github.com/abemedia/appcast/pkg/crypto/ed25519"
	"github.com/abemedia/appcast/source"
	"github.com/abemedia/appcast/target"
	"github.com/hashicorp/go-version"
	"github.com/imdario/mergo"
)

type Settings struct {
	InstallerArguments                string `yaml:"installerArguments"`
	MinimumSystemVersion              string `yaml:"minimumSystemVersion"`
	MinimumAutoupdateVersion          string `yaml:"minimumAutoupdateVersion"`
	IgnoreSkippedUpgradesBelowVersion string `yaml:"ignoreSkippedUpgradesBelowVersion"`
	CriticalUpdate                    bool   `yaml:"criticalUpdate"`
	CriticalUpdateBelowVersion        string `yaml:"criticalUpdateBelowVersion"`
}

type Rule struct {
	OS        OS     `yaml:"os"`
	Version   string `yaml:"version"`
	*Settings `yaml:",inline"`
}

type Config struct {
	Title       string
	Description string
	URL         string

	Source     *source.Source
	Target     target.Target
	FileName   string
	Version    string
	Prerelease bool
	DSAKey     *dsa.PrivateKey
	Ed25519Key ed25519.PrivateKey

	Settings []Rule
}

func getSettings(settings []Rule, v string, os OS) (*Settings, error) {
	opt := &Settings{}

	for _, s := range settings {
		if s.OS != Unknown && !IsOS(os, s.OS) {
			continue
		}

		if s.Version != "" && v == "" {
			continue
		}

		if s.Version != "" {
			constraints, err := version.NewConstraint(s.Version)
			if err != nil {
				return nil, err
			}

			if !constraints.Check(version.Must(version.NewVersion(v))) {
				continue
			}
		}

		if err := mergo.MergeWithOverwrite(opt, s.Settings); err != nil {
			return nil, err
		}
	}

	return opt, nil
}
