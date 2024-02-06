package config

import (
	"path"

	"github.com/kubri/kubri/integrations/sparkle"
	"github.com/kubri/kubri/pkg/crypto/dsa"
	"github.com/kubri/kubri/pkg/crypto/ed25519"
	"github.com/kubri/kubri/pkg/secret"
)

type sparkleConfig struct {
	Disabled    bool                  `yaml:"disabled,omitempty"`
	Folder      string                `yaml:"folder,omitempty"`
	Title       string                `yaml:"title,omitempty"`
	Description string                `yaml:"description,omitempty"`
	Filename    string                `yaml:"filename,omitempty"`
	DetectOS    map[sparkle.OS]string `yaml:"detect-os,omitempty"`
	Params      []struct {
		OS       sparkle.OS `yaml:"os,omitempty"      jsonschema:"type=string,enum=macos,enum=windows,enum=windows-x86,enum=windows-x64"` //nolint:lll
		Version  string     `yaml:"version,omitempty"`
		Settings *struct {
			InstallerArguments                string `yaml:"installer-arguments,omitempty"`
			MinimumSystemVersion              string `yaml:"minimum-system-version,omitempty"`
			MinimumAutoupdateVersion          string `yaml:"minimum-autoupdate-version,omitempty"`
			IgnoreSkippedUpgradesBelowVersion string `yaml:"ignore-skipped-upgrades-below-version,omitempty"`
			CriticalUpdate                    bool   `yaml:"critical-update,omitempty"`
			CriticalUpdateBelowVersion        string `yaml:"critical-update-below-version,omitempty"`
		} `yaml:",inline"`
	} `yaml:"params,omitempty"`
}

func getSparkle(c *config) (*sparkle.Config, error) {
	var dsaKey *dsa.PrivateKey
	if b, err := secret.Get("dsa_key"); err == nil {
		dsaKey, err = dsa.UnmarshalPrivateKey(b)
		if err != nil {
			return nil, err
		}
	}

	var edKey ed25519.PrivateKey
	if b, err := secret.Get("ed25519_key"); err == nil {
		edKey, err = ed25519.UnmarshalPrivateKey(b)
		if err != nil {
			return nil, err
		}
	}

	var detectOS func(string) sparkle.OS
	if len(c.Sparkle.DetectOS) > 0 {
		detectOS = func(s string) sparkle.OS {
			for k, v := range c.Sparkle.DetectOS {
				if ok, _ := path.Match(v, s); ok {
					return k
				}
			}
			return sparkle.Unknown
		}
	}

	params := make([]sparkle.Rule, len(c.Sparkle.Params))
	for i, p := range c.Sparkle.Params {
		params[i] = sparkle.Rule{
			OS:       p.OS,
			Version:  p.Version,
			Settings: (*sparkle.Settings)(p.Settings),
		}
	}

	return &sparkle.Config{
		Title:       fallback(c.Sparkle.Title, c.Title),
		Description: fallback(c.Sparkle.Description, c.Description),
		FileName:    fallback(c.Sparkle.Filename, "appcast.xml"),
		DSAKey:      dsaKey,
		Ed25519Key:  edKey,
		Settings:    params,
		DetectOS:    detectOS,

		Source:         c.source,
		Target:         c.target.Sub(fallback(c.Sparkle.Folder, "sparkle")),
		Version:        c.Version,
		Prerelease:     c.Prerelease,
		UploadPackages: c.UploadPackages,
	}, nil
}
