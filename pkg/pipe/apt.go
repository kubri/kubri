package pipe

import "github.com/abemedia/appcast/integrations/apt"

type aptConfig struct {
	Disabled bool   `yaml:"disabled"`
	Folder   string `yaml:"folder"`
}

func getApt(c *config) *apt.Config {
	dir := c.Apt.Folder
	if dir == "" {
		dir = "apt"
	}

	return &apt.Config{
		Source:     c.source,
		Target:     c.target.Sub(dir),
		Version:    c.Version,
		Prerelease: c.Prerelease,
	}
}
