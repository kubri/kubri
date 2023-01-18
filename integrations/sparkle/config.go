package sparkle

import (
	"github.com/abemedia/appcast/pkg/crypto/dsa"
	"github.com/abemedia/appcast/pkg/crypto/ed25519"
	"github.com/abemedia/appcast/pkg/version"
	"github.com/abemedia/appcast/source"
	"github.com/abemedia/appcast/target"
	"github.com/imdario/mergo"
)

type Settings struct {
	InstallerArguments                string
	MinimumSystemVersion              string
	MinimumAutoupdateVersion          string
	IgnoreSkippedUpgradesBelowVersion string
	CriticalUpdate                    bool
	CriticalUpdateBelowVersion        string
}

type Rule struct {
	OS      OS
	Version string
	*Settings
}

type Config struct {
	Title       string
	Description string
	URL         string
	DetectOS    func(string) OS

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

			if !constraints.Check(v) {
				continue
			}
		}

		if err := mergo.MergeWithOverwrite(opt, s.Settings); err != nil {
			return nil, err
		}
	}

	return opt, nil
}
