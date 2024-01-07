package sparkle

import (
	"dario.cat/mergo"
	"github.com/abemedia/appcast/pkg/crypto/dsa"
	"github.com/abemedia/appcast/pkg/crypto/ed25519"
	"github.com/abemedia/appcast/pkg/version"
	"github.com/abemedia/appcast/source"
	"github.com/abemedia/appcast/target"
)

type Config struct {
	Title       string
	Description string
	URL         string
	FileName    string
	DetectOS    func(string) OS
	DSAKey      *dsa.PrivateKey
	Ed25519Key  ed25519.PrivateKey
	Settings    []Rule

	Source         *source.Source
	Target         target.Target
	Version        string
	Prerelease     bool
	UploadPackages bool
}

type Rule struct {
	OS      OS
	Version string
	*Settings
}

type Settings struct {
	InstallerArguments                string
	MinimumSystemVersion              string
	MinimumAutoupdateVersion          string
	IgnoreSkippedUpgradesBelowVersion string
	CriticalUpdate                    bool
	CriticalUpdateBelowVersion        string
}

func getSettings(settings []Rule, v string, os OS) (*Settings, error) {
	opt := &Settings{}
	for _, s := range settings {
		if isOS(os, s.OS) && version.Check(s.Version, v) {
			if err := mergo.MergeWithOverwrite(opt, s.Settings); err != nil {
				return nil, err
			}
		}
	}
	return opt, nil
}
