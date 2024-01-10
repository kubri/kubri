package pipe

import "github.com/abemedia/appcast/integrations/appinstaller"

type appinstallerConfig struct {
	Disabled bool   `yaml:"disabled,omitempty"`
	Folder   string `yaml:"folder,omitempty"`
	OnLaunch *struct {
		HoursBetweenUpdateChecks int  `yaml:"hours-between-update-checks,omitempty"`
		ShowPrompt               bool `yaml:"show-prompt,omitempty"`
		UpdateBlocksActivation   bool `yaml:"update-blocks-activation,omitempty"`
	} `yaml:"on-launch,omitempty"`
	AutomaticBackgroundTask   bool `yaml:"automatic-background-task,omitempty"`
	ForceUpdateFromAnyVersion bool `yaml:"force-update-from-any-version,omitempty"`
}

func getAppinstaller(c *config) *appinstaller.Config {
	return &appinstaller.Config{
		OnLaunch:                  (*appinstaller.OnLaunchConfig)(c.Appinstaller.OnLaunch),
		AutomaticBackgroundTask:   c.Appinstaller.AutomaticBackgroundTask,
		ForceUpdateFromAnyVersion: c.Appinstaller.ForceUpdateFromAnyVersion,

		Source:         c.source,
		Target:         c.target.Sub(fallback(c.Appinstaller.Folder, "appinstaller")),
		Version:        c.Version,
		Prerelease:     c.Prerelease,
		UploadPackages: c.UploadPackages,
	}
}
