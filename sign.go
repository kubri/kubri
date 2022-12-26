package appcast

import (
	"encoding/base64"
	"fmt"
	"log"
	"strings"

	"github.com/abemedia/appcast/pkg/crypto/dsa"
	"github.com/abemedia/appcast/pkg/crypto/ed25519"
	"github.com/abemedia/appcast/pkg/os"
	"github.com/abemedia/appcast/source"
	"golang.org/x/mod/semver"
)

type SignOptions struct {
	Source     *source.Source
	Version    string
	DSAKey     *dsa.PrivateKey
	Ed25519Key ed25519.PrivateKey
}

// Sign signs update packages and uploads the signatures to the source.
func Sign(opt *SignOptions) error {
	if opt.Version != "" && strings.HasPrefix(opt.Version, semver.Canonical(opt.Version)) {
		release, err := opt.Source.GetRelease(opt.Version)
		if err != nil {
			return err
		}
		return signRelease(opt, release)
	}

	releases, err := opt.Source.ListReleases(&source.ListOptions{Constraint: opt.Version})
	if err != nil {
		return err
	}

	for _, release := range releases {
		err := signRelease(opt, release)
		if err != nil {
			return err
		}
	}

	return nil
}

func signRelease(opt *SignOptions, release *source.Release) error {
	if getAsset(release.Assets, "signatures.txt") != nil {
		return nil
	}

	sigs := signatures{}
	for _, asset := range release.Assets {
		algo, err := getAlgo(asset.Name)
		if err != nil {
			log.Printf("Skipping asset %s (%s): %s\n", asset.Name, release.Version, err)
			continue
		}

		log.Printf("Signing asset %s (%s)\n", asset.Name, release.Version)

		b, err := opt.Source.DownloadAsset(release.Version, asset.Name)
		if err != nil {
			return err
		}

		var sig []byte
		switch algo {
		case "ed25519":
			sig, err = ed25519.Sign(opt.Ed25519Key, b), nil
		case "dsa":
			sig, err = dsa.Sign(opt.DSAKey, b)
		}
		if err != nil {
			return err
		}

		sigs.Set(asset.Name, algo, base64.RawStdEncoding.EncodeToString(sig))
	}

	b, err := sigs.MarshalText()
	if err != nil {
		return err
	}

	return opt.Source.UploadAsset(release.Version, "signatures.txt", b)
}

func getAlgo(path string) (string, error) {
	switch os.Detect(path) {
	case os.MacOS:
		return "ed25519", nil
	case os.Windows, os.Windows64, os.Windows32:
		return "dsa", nil
	default:
		return "", fmt.Errorf("unsupported file extension")
	}
}
