package config

import (
	"cmp"

	"github.com/kubri/kubri/integrations/apk"
	"github.com/kubri/kubri/pkg/crypto/rsa"
	"github.com/kubri/kubri/pkg/secret"
)

type apkConfig struct {
	Disabled bool   `yaml:"disabled,omitempty"`
	Folder   string `yaml:"folder,omitempty"   validate:"omitempty,dirname"`
	KeyName  string `yaml:"key-name,omitempty"`
}

func getApk(c *config) (*apk.Config, error) {
	var rsaKey *rsa.PrivateKey
	if b, err := secret.Get("rsa_key"); err == nil {
		rsaKey, err = rsa.UnmarshalPrivateKey(b)
		if err != nil {
			return nil, err
		}
		if c.Apk.KeyName == "" {
			return nil, &Error{Errors: []string{"apk.key-name is required when rsa_key is set"}}
		}
	}

	return &apk.Config{
		Source:     c.source,
		Target:     c.target.Sub(cmp.Or(c.Apk.Folder, "apk")),
		Version:    c.Version,
		Prerelease: c.Prerelease,
		RSAKey:     rsaKey,
		KeyName:    c.Apk.KeyName,
	}, nil
}
