package pipe

import "github.com/abemedia/appcast/integrations/apt"

type aptConfig struct {
	Disabled bool   `yaml:"disabled"`
	Folder   string `yaml:"folder"`
}

func getApt(c *config) *apt.Config {
	return &apt.Config{
		Source:     c.source,
		Target:     c.target.Sub(fallback(c.Apt.Folder, "apt")),
		Version:    c.Version,
		Prerelease: c.Prerelease,
	}
}
