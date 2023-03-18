package apt

import (
	"crypto/md5"
	"crypto/sha1"
	"crypto/sha256"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"github.com/abemedia/appcast/integrations/apt/deb"
	"golang.org/x/mod/semver"
)

func release(p []*Package) (string, error) {
	dir, err := os.MkdirTemp("", "")
	if err != nil {
		return "", err
	}

	stable := make([]*Package, 0, len(p))
	for _, pkg := range p {
		if semver.Prerelease("v"+pkg.Version) == "" {
			stable = append(stable, pkg)
		}
	}
	if err = releaseSuite(stable, "stable", dir); err != nil {
		return "", err
	}

	// If not all packages are stable publish a separate `edge` dist.
	if len(p) > len(stable) {
		if err = releaseSuite(p, "edge", dir); err != nil {
			return "", err
		}
	}

	return dir, nil
}

func releaseSuite(p []*Package, suite, root string) error {
	dir := filepath.Join(root, "dists", suite)

	r := Releases{
		Suite:      suite,
		Codename:   suite,
		Date:       time.Now().UTC(),
		Components: "main",
	}

	byArch := map[string][]*Package{}
	for _, pkg := range p {
		byArch[pkg.Architecture] = append(byArch[pkg.Architecture], pkg)
	}

	{
		var as []string
		for a, pkgs := range byArch {
			if err := releaseArch(pkgs, suite, a, dir); err != nil {
				return err
			}
			as = append(as, a)
		}
		sort.Strings(as)
		r.Architectures = strings.Join(as, " ")
	}

	dirfs := os.DirFS(dir)
	err := fs.WalkDir(dirfs, ".", func(path string, d fs.DirEntry, err error) error {
		if err != nil || d.IsDir() {
			return err
		}

		b, err := fs.ReadFile(dirfs, path)
		if err != nil {
			return err
		}

		r.MD5Sum += fmt.Sprintf("\n%x %d %s", md5.Sum(b), len(b), path)
		r.SHA1 += fmt.Sprintf("\n%x %d %s", sha1.Sum(b), len(b), path)
		r.SHA256 += fmt.Sprintf("\n%x %d %s", sha256.Sum256(b), len(b), path)

		return nil
	})
	if err != nil {
		return err
	}

	if err := writeFile(filepath.Join(dir, "Release"), r); err != nil {
		return err
	}

	return writeFile(filepath.Join(dir, "InRelease"), r)
}

func releaseArch(p []*Package, suite, arch, root string) error {
	r := Release{
		Archive:      suite,
		Suite:        suite,
		Component:    "main",
		Architecture: arch,
	}

	dir := filepath.Join(root, "main", "binary-"+arch)
	if err := os.MkdirAll(dir, os.ModePerm); err != nil {
		return err
	}
	if err := writeFile(filepath.Join(dir, "Release"), r); err != nil {
		return err
	}
	if err := writeFile(filepath.Join(dir, "Packages"), p); err != nil {
		return err
	}
	return writeFile(filepath.Join(dir, "Packages.gz"), p)
}

func writeFile(path string, v any) error {
	f, err := os.Create(path)
	if err != nil {
		return err
	}
	w, err := compress(filepath.Ext(path))(f)
	if err != nil {
		return err
	}
	if err = deb.NewEncoder(w).Encode(v); err != nil {
		return err
	}
	if err = w.Close(); err != nil {
		return err
	}
	return f.Close()
}
