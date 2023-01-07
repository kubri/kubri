package apt

import (
	"context"
	"crypto/md5"
	"crypto/sha1"
	"crypto/sha256"
	"fmt"
	"io/fs"
	"os"
	"path"
	"strings"

	"github.com/abemedia/appcast/integrations/apt/deb"
	"github.com/abemedia/appcast/source"
	"github.com/abemedia/appcast/target"
)

type Config struct {
	Source     *source.Source
	Version    string
	Prerelease bool
	Target     target.Target
}

func Build(ctx context.Context, c *Config) error {
	releases, err := c.Source.ListReleases(ctx, &source.ListOptions{
		Version:    c.Version,
		Prerelease: c.Prerelease,
	})
	if err != nil {
		return err
	}

	cached := read(ctx, c)
	var isNew bool
	var items []*Package
	for _, release := range releases {
		for _, asset := range release.Assets {
			if path.Ext(asset.Name) != ".deb" {
				continue
			}
			if item, ok := cached[asset.Name]; ok {
				items = append(items, item)
				continue
			}
			isNew = true
			item, err := getPackage(ctx, c, release.Version, asset.Name)
			if err != nil {
				return err
			}
			items = append(items, item)
		}
	}
	if !isNew {
		return nil // No new packages.
	}

	dir, err := release(items)
	if err != nil {
		return err
	}
	defer os.RemoveAll(dir)

	files := os.DirFS(dir)
	return fs.WalkDir(files, ".", func(path string, d fs.DirEntry, err error) error {
		if err != nil || d.IsDir() {
			return err
		}
		w, err := c.Target.NewWriter(ctx, path)
		if err != nil {
			return err
		}
		b, err := fs.ReadFile(files, path)
		if err != nil {
			return err
		}
		if _, err = w.Write(b); err != nil {
			return err
		}
		return w.Close()
	})
}

func getPackage(ctx context.Context, c *Config, version, name string) (*Package, error) {
	b, err := c.Source.DownloadAsset(ctx, version, name)
	if err != nil {
		return nil, err
	}

	p, err := getControl(b)
	if err != nil {
		return nil, err
	}
	p.Size = len(b)
	p.Filename = path.Join("pool/main", p.Package[0:1], p.Package, name)
	p.MD5sum = md5.Sum(b)
	p.SHA1 = sha1.Sum(b)
	p.SHA256 = sha256.Sum256(b)

	w, err := c.Target.NewWriter(ctx, p.Filename)
	if err != nil {
		return nil, err
	}
	if _, err = w.Write(b); err != nil {
		return nil, err
	}
	if err = w.Close(); err != nil {
		return nil, err
	}

	return p, nil
}

func read(ctx context.Context, c *Config) map[string]*Package {
	res := map[string]*Package{}

	dist := "edge"
	rd, err := c.Target.NewReader(ctx, "dists/edge/Release")
	if err != nil {
		dist = "stable"
		rd, err = c.Target.NewReader(ctx, "dists/stable/Release")
		if err != nil {
			return res
		}
	}
	defer rd.Close()

	r := &Releases{}
	if err = deb.NewDecoder(rd).Decode(r); err != nil {
		return res
	}

	var pkgs []*Package
	for _, arch := range strings.Split(r.Architectures, " ") {
		path := fmt.Sprintf("dists/%s/main/binary-%s/Packages", dist, arch)
		rd, err = c.Target.NewReader(ctx, path)
		if err != nil {
			return res
		}

		var p []*Package
		if err = deb.NewDecoder(rd).Decode(&p); err != nil {
			return res
		}
		pkgs = append(pkgs, p...)
	}

	for _, p := range pkgs {
		res[path.Base(p.Filename)] = p
	}

	return res
}
