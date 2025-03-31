package config

import (
	"cmp"

	"github.com/kubri/kubri/integrations/yum"
	"github.com/kubri/kubri/pkg/crypto/pgp"
	"github.com/kubri/kubri/pkg/secret"
)

type yumConfig struct {
	Disabled bool   `yaml:"disabled,omitempty"`
	Folder   string `yaml:"folder,omitempty"   validate:"omitempty,dirname"`
}

func getYum(c *config) (*yum.Config, error) {
	var pgpKey *pgp.PrivateKey
	if b, err := secret.Get("pgp_key"); err == nil {
		pgpKey, err = pgp.UnmarshalPrivateKey(b)
		if err != nil {
			return nil, err
		}
	}

	return &yum.Config{
		Source:     c.source,
		Target:     c.target.Sub(cmp.Or(c.Yum.Folder, "yum")),
		Version:    c.Version,
		Prerelease: c.Prerelease,
		PGPKey:     pgpKey,
	}, nil
}
