package config

import (
	"cmp"

	"github.com/kubri/kubri/integrations/arch"
	"github.com/kubri/kubri/pkg/crypto/pgp"
	"github.com/kubri/kubri/pkg/secret"
)

type archConfig struct {
	Disabled bool   `yaml:"disabled,omitempty"`
	Folder   string `yaml:"folder,omitempty"   validate:"omitempty,dirname"`
	RepoName string `yaml:"repo-name"          validate:"required,slug"`
}

func getArch(c *config) (*arch.Config, error) {
	var pgpKey *pgp.PrivateKey
	if b, err := secret.Get("pgp_key"); err == nil {
		pgpKey, err = pgp.UnmarshalPrivateKey(b)
		if err != nil {
			return nil, err
		}
	}

	return &arch.Config{
		RepoName:   c.Arch.RepoName,
		Source:     c.source,
		Target:     c.target.Sub(cmp.Or(c.Arch.Folder, "arch")),
		Version:    c.Version,
		Prerelease: c.Prerelease,
		PGPKey:     pgpKey,
	}, nil
}
