package pipe

import "github.com/abemedia/appcast/integrations/appinstaller"

type appinstallerConfig struct {
	Disabled                  bool   `yaml:"disabled"`
	Folder                    string `yaml:"folder"`
	HoursBetweenUpdateChecks  int    `yaml:"hours_between_update_checks"`
	UpdateBlocksActivation    bool   `yaml:"update_blocks_activation"`
	ShowPrompt                bool   `yaml:"show_prompt"`
	AutomaticBackgroundTask   bool   `yaml:"automatic_background_task"`
	ForceUpdateFromAnyVersion bool   `yaml:"force_update_from_any_version"`
}

func getAppinstaller(c *config) *appinstaller.Config {
	return &appinstaller.Config{
		HoursBetweenUpdateChecks:  c.Appinstaller.HoursBetweenUpdateChecks,
		UpdateBlocksActivation:    c.Appinstaller.UpdateBlocksActivation,
		ShowPrompt:                c.Appinstaller.ShowPrompt,
		AutomaticBackgroundTask:   c.Appinstaller.AutomaticBackgroundTask,
		ForceUpdateFromAnyVersion: c.Appinstaller.ForceUpdateFromAnyVersion,

		Source:         c.source,
		Target:         c.target.Sub(fallback(c.Appinstaller.Folder, "appinstaller")),
		Version:        c.Version,
		Prerelease:     c.Prerelease,
		UploadPackages: c.UploadPackages,
	}
}
