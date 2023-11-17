package pipe

import (
	"github.com/abemedia/appcast/integrations/apt"
	"github.com/abemedia/appcast/pkg/crypto/pgp"
	"github.com/abemedia/appcast/pkg/secret"
)

type aptConfig struct {
	Disabled bool   `yaml:"disabled"`
	Folder   string `yaml:"folder"`
}

func getApt(c *config) (*apt.Config, error) {
	var pgpKey *pgp.PrivateKey
	if b, err := secret.Get("pgp_key"); err == nil {
		pgpKey, err = pgp.UnmarshalPrivateKey(b)
		if err != nil {
			return nil, err
		}
	}

	return &apt.Config{
		Source:     c.source,
		Target:     c.target.Sub(fallback(c.Apt.Folder, "apt")),
		Version:    c.Version,
		Prerelease: c.Prerelease,
		PGPKey:     pgpKey,
	}, nil
}
