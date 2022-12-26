package appcast

import (
	"github.com/abemedia/appcast/pkg/os"
	"github.com/abemedia/appcast/source"
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
	OS        os.OS  `yaml:"os"`
	Version   string `yaml:"version"`
	*Settings `yaml:",inline"`
}

type Config struct {
	Title       string `yaml:"title"`
	Description string `yaml:"description"`
	URL         string `yaml:"url"`

	Source     *source.Source `yaml:"source"`
	Prerelease bool           `yaml:"prerelease"`
	RewriteURL RewriteFunc    `yaml:"rewriteUrl"`

	Settings []Rule `yaml:"settings"`
}

func getSettings(settings []Rule, v string, o os.OS) (*Settings, error) {
	opt := &Settings{}

	for _, s := range settings {
		if s.OS != os.Unknown && !os.Is(o, s.OS) {
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
