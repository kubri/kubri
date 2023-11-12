package pipe

import "github.com/abemedia/appcast/integrations/yum"

type yumConfig struct {
	Disabled bool   `yaml:"disabled"`
	Folder   string `yaml:"folder"`
}

func getYum(c *config) *yum.Config {
	return &yum.Config{
		Source:     c.source,
		Target:     c.target.Sub(fallback(c.Yum.Folder, "yum")),
		Version:    c.Version,
		Prerelease: c.Prerelease,
	}
}
