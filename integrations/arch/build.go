package arch

import (
	"context"
	"os"
	"strings"

	"github.com/kubri/kubri/pkg/crypto/pgp"
	"github.com/kubri/kubri/source"
	"github.com/kubri/kubri/target"
)

type Config struct {
	RepoName   string
	Source     *source.Source
	Version    string
	Prerelease bool
	Target     target.Target
	PGPKey     *pgp.PrivateKey
}

func Build(ctx context.Context, c *Config) error {
	r, err := openRepo(ctx, c.Target, c.RepoName, c.PGPKey)
	if err != nil {
		return err
	}
	defer os.RemoveAll(r.dir)

	version := c.Version
	if v := getLatest(r.packages); v != "" {
		version += ",>v" + v
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

	var hasNew bool

	for _, rel := range releases {
		for _, asset := range rel.Assets {
			if !isValidPackage(asset.Name) {
				continue
			}
			data, err := c.Source.DownloadAsset(ctx, rel.Version, asset.Name)
			if err != nil {
				return err
			}
			if err := r.Add(asset.Name, data); err != nil {
				return err
			}
			hasNew = true
		}
	}

	if !hasNew {
		return nil
	}

	if err := r.Write(); err != nil {
		return err
	}

	if err := target.CopyFS(ctx, c.Target, os.DirFS(r.dir)); err != nil {
		return err
	}

	return nil
}

func isValidPackage(filename string) bool {
	validExts := []string{".pkg.tar.zst", ".pkg.tar.gz", ".pkg.tar.xz", ".pkg.tar.bz2"}
	for _, ext := range validExts {
		if strings.HasSuffix(filename, ext) {
			return true
		}
	}
	return false
}

func getLatest(repo map[string]map[string]map[string]*Package) string {
	var latest string
	for _, pkgMap := range repo {
		for _, versions := range pkgMap {
			for _, pkg := range versions {
				if latest == "" || compareVersions(pkg.Version, latest) > 0 {
					latest = pkg.Version
				}
			}
		}
	}
	return stripVersion(latest)
}

func stripVersion(version string) string {
	if colon := strings.IndexByte(version, ':'); colon != -1 {
		version = version[colon+1:]
	}
	if dash := strings.LastIndexByte(version, '-'); dash != -1 {
		version = version[:dash]
	}
	return strings.ReplaceAll(version, "_", "-")
}
