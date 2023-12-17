package pipe

import "github.com/abemedia/appcast/integrations/appinstaller"

type appinstallerConfig struct {
	Disabled bool   `yaml:"disabled"`
	Folder   string `yaml:"folder"`
	OnLaunch *struct {
		HoursBetweenUpdateChecks int  `yaml:"hours-between-update-checks"`
		ShowPrompt               bool `yaml:"show-prompt"`
		UpdateBlocksActivation   bool `yaml:"update-blocks-activation"`
	} `yaml:"on-launch"`
	AutomaticBackgroundTask   bool `yaml:"automatic-background-task"`
	ForceUpdateFromAnyVersion bool `yaml:"force-update-from-any-version"`
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
