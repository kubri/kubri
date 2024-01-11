package pipe

import (
	"fmt"

	"github.com/abemedia/appcast/integrations/apt"
	"github.com/abemedia/appcast/pkg/crypto/pgp"
	"github.com/abemedia/appcast/pkg/secret"
)

type aptConfig struct {
	Disabled bool     `yaml:"disabled,omitempty"`
	Folder   string   `yaml:"folder,omitempty"   validate:"omitempty,dirname"`
	Compress []string `yaml:"compress,omitempty" validate:"dive,oneof=none gzip bzip2 xz lzma lz4 zstd" jsonschema:"enum=none,enum=gzip,enum=bzip2,enum=xz,enum=lzma,enum=lz4,enum=zstd"` //nolint:lll
}

func getApt(c *config) (*apt.Config, error) {
	var pgpKey *pgp.PrivateKey
	if b, err := secret.Get("pgp_key"); err == nil {
		pgpKey, err = pgp.UnmarshalPrivateKey(b)
		if err != nil {
			return nil, err
		}
	}

	var algos apt.CompressionAlgo
	for _, name := range c.Apt.Compress {
		switch name {
		case "none":
			algos |= apt.NoCompression
		case "gzip":
			algos |= apt.GZIP
		case "bzip2":
			algos |= apt.BZIP2
		case "xz":
			algos |= apt.XZ
		case "lzma":
			algos |= apt.LZMA
		case "lz4":
			algos |= apt.LZ4
		case "zstd":
			algos |= apt.ZSTD
		default:
			return nil, fmt.Errorf("unknown compression algorithm: %s", name)
		}
	}

	return &apt.Config{
		Source:     c.source,
		Target:     c.target.Sub(fallback(c.Apt.Folder, "apt")),
		Version:    c.Version,
		Prerelease: c.Prerelease,
		PGPKey:     pgpKey,
		Compress:   algos,
	}, nil
}
