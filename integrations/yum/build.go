package yum

import (
	"context"
	"io/fs"
	"log"
	"os"
	"path"
	"strings"
	"unsafe"

	"github.com/kubri/kubri/pkg/crypto/pgp"
	"github.com/kubri/kubri/source"
	"github.com/kubri/kubri/target"
)

type Config struct {
	Source     *source.Source
	Version    string
	Prerelease bool
	Target     target.Target
	PGPKey     *pgp.PrivateKey
}

// Build creates or updates a YUM repository.
//
//nolint:funlen,gocognit
func Build(ctx context.Context, c *Config) error {
	repo, err := openRepo(ctx, c.Target)
	if err != nil {
		return err
	}
	defer os.RemoveAll(repo.dir)

	version := c.Version
	if v := getVersionConstraint(repo.primary.Package); v != "" {
		version += "," + v
	}

	releases, err := c.Source.ListReleases(ctx, &source.ListOptions{
		Version:    version,
		Prerelease: c.Prerelease,
	})
	if err == source.ErrNoReleaseFound {
		return nil
	}
	if err != nil {
		return err
	}

	var hasReleases bool
	for _, release := range releases {
		for _, asset := range release.Assets {
			if path.Ext(asset.Name) != ".rpm" {
				continue
			}
			b, err := c.Source.DownloadAsset(ctx, release.Version, asset.Name)
			if err != nil {
				return err
			}
			if err = repo.Add(b); err != nil {
				return err
			}
			hasReleases = true
		}
	}
	if !hasReleases {
		return nil
	}

	if err = repo.Write(c.PGPKey); err != nil {
		return err
	}

	files := os.DirFS(repo.dir)
	err = fs.WalkDir(files, ".", func(path string, d fs.DirEntry, err error) error {
		if err != nil || d.IsDir() {
			return err
		}
		b, err := fs.ReadFile(files, path)
		if err != nil {
			return err
		}
		w, err := c.Target.NewWriter(ctx, path)
		if err != nil {
			return err
		}
		if _, err = w.Write(b); err != nil {
			return err
		}
		return w.Close()
	})
	if err != nil {
		return err
	}

	for _, path := range repo.files {
		if err = c.Target.Remove(ctx, path); err != nil {
			log.Printf("Failed to delete %s: %s", path, err)
		}
	}

	return nil
}

func getVersionConstraint(pkgs []Package) string {
	if len(pkgs) == 0 {
		return ""
	}

	v := make([]byte, 0, len(pkgs)*len("!=0.0.0,"))
	for _, p := range pkgs {
		v = append(v, '!', '=')
		v = append(v, strings.Replace(p.Version.Ver, "~", "-", 1)...)
		v = append(v, ',')
	}

	return unsafe.String(unsafe.SliceData(v), len(v)-1)
}
