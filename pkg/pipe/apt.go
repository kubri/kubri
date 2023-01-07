package pipe

import "github.com/abemedia/appcast/integrations/apt"

type aptConfig struct {
	Disabled bool `yaml:"disabled"`
}

func getApt(c *config) *apt.Config {
	tgt := c.target
	if !c.Target.Flat {
		tgt = tgt.Sub("apt")
	}

	return &apt.Config{
		Source:     c.source,
		Target:     tgt,
		Version:    c.Source.Version,
		Prerelease: c.Source.Prerelease,
	}
}
