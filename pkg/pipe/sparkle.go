package pipe

import (
	"path"

	"github.com/abemedia/appcast/integrations/sparkle"
	"github.com/abemedia/appcast/pkg/crypto/dsa"
	"github.com/abemedia/appcast/pkg/crypto/ed25519"
	"github.com/abemedia/appcast/pkg/secret"
)

type sparkleConfig struct {
	Disabled    bool                 `yaml:"disabled"`
	Folder      string               `yaml:"folder"`
	Title       string               `yaml:"title"`
	Description string               `yaml:"description"`
	Filename    string               `yaml:"filename"`
	DetectOS    map[sparkleOS]string `yaml:"detect_os"`
	Params      []struct {
		OS       sparkleOS `yaml:"os"`
		Version  string    `yaml:"version"`
		Settings *struct {
			InstallerArguments                string `yaml:"installer_arguments"`
			MinimumSystemVersion              string `yaml:"minimum_system_version"`
			MinimumAutoupdateVersion          string `yaml:"minimum_autoupdate_version"`
			IgnoreSkippedUpgradesBelowVersion string `yaml:"ignore_skipped_upgrades_below_version"`
			CriticalUpdate                    bool   `yaml:"critical_update"`
			CriticalUpdateBelowVersion        string `yaml:"critical_update_below_version"`
		} `yaml:",inline"`
	} `yaml:"params"`
}

type sparkleOS sparkle.OS

func (os *sparkleOS) UnmarshalText(text []byte) error {
	switch string(text) {
	default:
		*os = sparkleOS(sparkle.Unknown)
	case "macos":
		*os = sparkleOS(sparkle.MacOS)
	case "windows":
		*os = sparkleOS(sparkle.Windows)
	case "windows_x64":
		*os = sparkleOS(sparkle.Windows64)
	case "windows_x86":
		*os = sparkleOS(sparkle.Windows32)
	}
	return nil
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
					return sparkle.OS(k)
				}
			}
			return sparkle.Unknown
		}
	}

	params := make([]sparkle.Rule, len(c.Sparkle.Params))
	for i, p := range c.Sparkle.Params {
		params[i] = sparkle.Rule{
			OS:       sparkle.OS(p.OS),
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
