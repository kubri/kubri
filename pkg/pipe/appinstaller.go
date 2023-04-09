package pipe

import "github.com/abemedia/appcast/integrations/appinstaller"

type appinstallerConfig struct {
	Disabled                  bool   `yaml:"disabled"`
	Folder                    string `yaml:"folder"`
	HoursBetweenUpdateChecks  int    `yaml:"hours-between-update-checks"`
	UpdateBlocksActivation    bool   `yaml:"update-blocks-activation"`
	ShowPrompt                bool   `yaml:"show-prompt"`
	AutomaticBackgroundTask   bool   `yaml:"automatic-background-task"`
	ForceUpdateFromAnyVersion bool   `yaml:"force-update-from-any-version"`
}

func getAppinstaller(c *config) *appinstaller.Config {
	if c.Appinstaller.Disabled {
		return nil
	}

	dir := c.Appinstaller.Folder
	if dir == "" {
		dir = "appinstaller"
	}

	return &appinstaller.Config{
		HoursBetweenUpdateChecks:  c.Appinstaller.HoursBetweenUpdateChecks,
		UpdateBlocksActivation:    c.Appinstaller.UpdateBlocksActivation,
		ShowPrompt:                c.Appinstaller.ShowPrompt,
		AutomaticBackgroundTask:   c.Appinstaller.AutomaticBackgroundTask,
		ForceUpdateFromAnyVersion: c.Appinstaller.ForceUpdateFromAnyVersion,

		Source:         c.source,
		Target:         c.target.Sub(dir),
		Version:        c.Version,
		Prerelease:     c.Prerelease,
		UploadPackages: c.UploadPackages,
	}
}
