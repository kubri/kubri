package apt

import (
	"crypto/md5"
	"crypto/sha1"
	"crypto/sha256"
	"io/fs"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"golang.org/x/mod/semver"

	"github.com/kubri/kubri/integrations/apt/deb"
	"github.com/kubri/kubri/pkg/crypto/pgp"
)

func release(key *pgp.PrivateKey, algos CompressionAlgo, p []*Package) (string, error) {
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
	if err = releaseSuite(key, algos, stable, "stable", dir); err != nil {
		return "", err
	}

	// If not all packages are stable publish a separate `edge` dist.
	if len(p) > len(stable) {
		if err = releaseSuite(key, algos, p, "edge", dir); err != nil {
			return "", err
		}
	}

	if key != nil {
		b, err := pgp.MarshalPublicKey(pgp.Public(key))
		if err != nil {
			return "", err
		}
		if err = os.WriteFile(filepath.Join(dir, "key.asc"), b, 0o600); err != nil {
			return "", err
		}
	}

	return dir, nil
}

func releaseSuite(key *pgp.PrivateKey, algos CompressionAlgo, p []*Package, suite, root string) error {
	dir := filepath.Join(root, "dists", suite)

	r := Releases{
		Suite:      suite,
		Codename:   suite,
		Components: "main",
	}

	byArch := map[string][]*Package{}
	for _, pkg := range p {
		byArch[pkg.Architecture] = append(byArch[pkg.Architecture], pkg)
	}

	{
		var as []string
		for a, pkgs := range byArch {
			if err := releaseArch(algos, pkgs, suite, a, dir); err != nil {
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

		// TODO: Do we want this? What is perf impact vs benefit?
		r.MD5Sum = append(r.MD5Sum, ChecksumFile[[16]byte]{md5.Sum(b), len(b), path})
		r.SHA1 = append(r.SHA1, ChecksumFile[[20]byte]{sha1.Sum(b), len(b), path})
		r.SHA256 = append(r.SHA256, ChecksumFile[[32]byte]{sha256.Sum256(b), len(b), path})

		return nil
	})
	if err != nil {
		return err
	}

	return writeRelease(dir, r, key)
}

func releaseArch(algos CompressionAlgo, p []*Package, suite, arch, root string) error {
	r := Release{
		Archive:      suite,
		Suite:        suite,
		Component:    "main",
		Architecture: arch,
	}

	dir := filepath.Join(root, "main", "binary-"+arch)
	if err := os.MkdirAll(dir, 0o750); err != nil {
		return err
	}
	if err := writeFile(filepath.Join(dir, "Release"), r); err != nil {
		return err
	}
	return writePackages(filepath.Join(dir, "Packages"), p, algos)
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

func writePackages(path string, v []*Package, algos CompressionAlgo) error {
	b, err := deb.Marshal(v)
	if err != nil {
		return err
	}

	for _, ext := range compressionExtensions(algos) {
		f, err := os.Create(path + ext)
		if err != nil {
			return err
		}
		w, err := compress(ext)(f)
		if err != nil {
			return err
		}
		if _, err = w.Write(b); err != nil {
			return err
		}
		if err = w.Close(); err != nil {
			return err
		}
		if err = f.Close(); err != nil {
			return err
		}
	}

	return nil
}

var timeNow = time.Now //nolint:gochecknoglobals

func writeRelease(dir string, r Releases, key *pgp.PrivateKey) error {
	r.Date = timeNow().UTC()

	b, err := deb.Marshal(r)
	if err != nil {
		return err
	}

	if err = os.WriteFile(filepath.Join(dir, "Release"), b, 0o600); err != nil {
		return err
	}

	if key != nil {
		sig, err := pgp.Sign(key, b)
		if err != nil {
			return err
		}
		if err = os.WriteFile(filepath.Join(dir, "Release.gpg"), sig, 0o600); err != nil {
			return err
		}
		b, err = pgp.SignText(key, b)
		if err != nil {
			return err
		}
	}

	return os.WriteFile(filepath.Join(dir, "InRelease"), b, 0o600)
}
