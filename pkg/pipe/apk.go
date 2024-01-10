package pipe

import (
	"errors"

	"github.com/abemedia/appcast/integrations/apk"
	"github.com/abemedia/appcast/pkg/crypto/rsa"
	"github.com/abemedia/appcast/pkg/secret"
)

type apkConfig struct {
	Disabled bool   `yaml:"disabled,omitempty"`
	Folder   string `yaml:"folder,omitempty"`
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
			return nil, errors.New("missing key name")
		}
	}

	return &apk.Config{
		Source:     c.source,
		Target:     c.target.Sub(fallback(c.Apk.Folder, "apk")),
		Version:    c.Version,
		Prerelease: c.Prerelease,
		RSAKey:     rsaKey,
		KeyName:    c.Apk.KeyName,
	}, nil
}
