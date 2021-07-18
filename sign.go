package appcast

import (
	"encoding/base64"
	"encoding/pem"
	"fmt"
	"os"

	"github.com/abemedia/appcast/pkg/crypto/dsa"
	"github.com/abemedia/appcast/pkg/crypto/ed25519"
)

func Sign(c *Config) error {
	if c.Source == nil {
		return fmt.Errorf("missing source")
	}

	sign := makeSigner(c)

	releases, err := c.Source.Releases()
	if err != nil {
		return err
	}

	for _, release := range releases {
		s := getAsset(release.Assets, "signatures.txt")
		if s != nil {
			continue
		}

		sigs := signatures{}
		for _, asset := range release.Assets {
			os := detectOS(c, asset.Name)
			if os == Unknown {
				continue
			}

			b, err := c.Source.DownloadAsset(release.Version, asset.Name)
			if err != nil {
				return err
			}

			algorithm := getAlgorithm(os)

			sig, err := sign(algorithm, b)
			if err != nil {
				return err
			}

			sigs.Set(asset.Name, algorithm, base64.RawStdEncoding.EncodeToString(sig))
		}

		b, err := sigs.MarshalText()
		if err != nil {
			return err
		}

		err = c.Source.UploadAsset(release.Version, "signatures.txt", b)
		if err != nil {
			return err
		}
	}

	return nil
}

func makeSigner(c *Config) func(algo string, b []byte) ([]byte, error) {
	var dsaKey *dsa.PrivateKey
	var edKey ed25519.PrivateKey

	return func(algo string, b []byte) ([]byte, error) {
		switch algo {
		case "ed25519":
			if edKey == nil {
				der, err := readPEMFile(c.EdKey)
				if err != nil {
					return nil, fmt.Errorf("failed to read %s key: %w", algo, err)
				}
				edKey, err = ed25519.UnmarshalPrivateKey(der)
				if err != nil {
					return nil, err
				}
			}
			return ed25519.Sign(edKey, b), nil

		case "dsa":
			if dsaKey == nil {
				der, err := readPEMFile(c.DSAKey)
				if err != nil {
					return nil, fmt.Errorf("failed to read %s key: %w", algo, err)
				}
				dsaKey, err = dsa.UnmarshalPrivateKey(der)
				if err != nil {
					return nil, err
				}
			}
			return dsa.Sign(dsaKey, b)

		default:
			return nil, fmt.Errorf("unknown algorithm: %s", algo)
		}
	}
}

func readPEMFile(path string) ([]byte, error) {
	if path == "" {
		return nil, fmt.Errorf("no key provided")
	}

	b, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	p, _ := pem.Decode(b)
	if p == nil {
		return nil, fmt.Errorf("no pem block found")
	}

	return p.Bytes, nil
}

func getAlgorithm(os OS) string {
	switch os {
	case MacOS:
		return "ed25519"
	case Windows64, Windows32, Windows:
		return "dsa"
	default:
		panic("unknown OS: " + os.String())
	}
}
