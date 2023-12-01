package apk

import (
	"context"
	"io/fs"
	"os"
	"path"
	"unsafe"

	"github.com/abemedia/appcast/pkg/crypto/rsa"
	"github.com/abemedia/appcast/source"
	"github.com/abemedia/appcast/target"
	"gitlab.alpinelinux.org/alpine/go/repository"
)

type Config struct {
	Source     *source.Source
	Version    string
	Prerelease bool
	Target     target.Target
	RSAKey     *rsa.PrivateKey
	KeyName    string
}

//nolint:funlen
func Build(ctx context.Context, c *Config) error {
	repo, err := openRepo(ctx, c.Target)
	if err != nil {
		return err
	}
	defer os.RemoveAll(repo.dir)

	version := c.Version
	if v := getVersionConstraint(repo.repos); v != "" {
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
			if path.Ext(asset.Name) != ".apk" {
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

	if err = repo.Write(c.RSAKey, c.KeyName); err != nil {
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

	return nil
}

func getVersionConstraint(repo map[string]*repository.ApkIndex) string {
	if len(repo) == 0 {
		return ""
	}

	v := make([]byte, 0)
	for _, r := range repo {
		for _, p := range r.Packages {
			v = append(v, '!', '=')
			v = append(v, p.Version...)
			v = append(v, ',')
		}
	}

	return unsafe.String(unsafe.SliceData(v), len(v)-1)
}
