package apt

import (
	"context"
	"crypto/md5"
	"crypto/sha1"
	"crypto/sha256"
	"fmt"
	"os"
	"path"
	"strings"
	"unsafe"

	"github.com/kubri/kubri/integrations/apt/deb"
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
	Compress   CompressionAlgo
}

func Build(ctx context.Context, c *Config) error {
	pkgs := read(ctx, c)

	version := c.Version
	if v := getVersionConstraint(pkgs); v != "" {
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

	p, err := getPackages(ctx, c, releases)
	if err != nil || p == nil {
		return err
	}
	pkgs = append(p, pkgs...)

	dir, err := release(c.PGPKey, c.Compress, pkgs)
	if err != nil {
		return err
	}
	defer os.RemoveAll(dir)

	return target.CopyFS(ctx, c.Target, os.DirFS(dir))
}

func read(ctx context.Context, c *Config) []*Package {
	dist := "edge"
	rd, err := c.Target.NewReader(ctx, "dists/edge/Release")
	if err != nil {
		dist = "stable"
		rd, err = c.Target.NewReader(ctx, "dists/stable/Release")
		if err != nil {
			return nil
		}
	}
	defer rd.Close()

	r := &Releases{}
	if err = deb.NewDecoder(rd).Decode(r); err != nil {
		return nil
	}

	var pkgs []*Package
	for arch := range strings.SplitSeq(r.Architectures, " ") {
		path := fmt.Sprintf("dists/%s/main/binary-%s/Packages", dist, arch)
		rd, err = c.Target.NewReader(ctx, path)
		if err != nil {
			return nil
		}
		defer rd.Close()

		var p []*Package
		if err = deb.NewDecoder(rd).Decode(&p); err != nil {
			return nil
		}
		pkgs = append(pkgs, p...)
	}

	return pkgs
}

func getVersionConstraint(pkgs []*Package) string {
	if len(pkgs) == 0 {
		return ""
	}

	v := make([]byte, 0, len(pkgs)*len("!=0.0.0,"))
	for _, p := range pkgs {
		v = append(v, '!', '=')
		v = append(v, strings.Replace(p.Version, "~", "-", 1)...)
		v = append(v, ',')
	}

	return unsafe.String(unsafe.SliceData(v), len(v)-1)
}

func getPackages(ctx context.Context, c *Config, releases []*source.Release) ([]*Package, error) {
	var pkgs []*Package
	for _, release := range releases {
		for _, asset := range release.Assets {
			if path.Ext(asset.Name) != ".deb" {
				continue
			}
			p, err := getPackage(ctx, c, release.Version, asset.Name)
			if err != nil {
				return nil, err
			}
			pkgs = append(pkgs, p)
		}
	}
	return pkgs, nil
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
	p.Filename = "pool/main/" + p.Package[0:1] + "/" + p.Package + "/" +
		p.Package + "_" + p.Version + "_" + p.Architecture + ".deb"
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
