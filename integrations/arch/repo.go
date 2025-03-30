package arch

import (
	"archive/tar"
	"bytes"
	"context"
	"crypto/sha256"
	"errors"
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/klauspost/compress/zstd"

	"github.com/kubri/kubri/integrations/arch/desc"
	"github.com/kubri/kubri/integrations/arch/pkginfo"
	"github.com/kubri/kubri/pkg/crypto/pgp"
	"github.com/kubri/kubri/target"
)

type repo struct {
	dir      string
	name     string
	pgpKey   *pgp.PrivateKey
	packages map[string]map[string]map[string]*Package // arch -> pkgName -> versions
}

func openRepo(ctx context.Context, t target.Target, repoName string, pgpKey *pgp.PrivateKey) (*repo, error) {
	dir, err := os.MkdirTemp("", "archrepo-")
	if err != nil {
		return nil, err
	}

	r := &repo{
		dir:      dir,
		name:     repoName,
		pgpKey:   pgpKey,
		packages: map[string]map[string]map[string]*Package{},
	}

	// TODO: Implement `Walk` on targets so we can iterate over all folders and
	// avoid hard-coding the archs.
	archs := []string{
		"x86_64", "any",
		"aarch64", "armv7h", // https://archlinuxarm.org/packages
		"powerpc64le", "powerpc64", "powerpc", "riscv64", // https://archlinuxpower.org/
		"i686", "pentium4", // https://archlinux32.org/architecture/
	}

	for _, arch := range archs {
		dbPath := filepath.Join(arch, repoName+".db")
		db, err := t.NewReader(ctx, dbPath)
		if err != nil {
			if errors.Is(err, fs.ErrNotExist) {
				continue
			}
			return nil, err
		}
		defer db.Close()

		if err := r.readDB(db, arch); err != nil {
			return r, fmt.Errorf("failed to parse %s: %w", dbPath, err)
		}
	}

	return r, nil
}

func (r *repo) readDB(dbData io.ReadCloser, arch string) error {
	zr, err := zstd.NewReader(dbData)
	if err != nil {
		return fmt.Errorf("failed to decompress: %w", err)
	}
	defer zr.Close()

	tarReader := tar.NewReader(zr)

	for {
		hdr, err := tarReader.Next()
		if errors.Is(err, io.EOF) {
			break
		}
		if err != nil {
			return err
		}
		if hdr.Typeflag != tar.TypeReg || !strings.HasSuffix(hdr.Name, "/desc") {
			continue
		}

		var pkg Package
		if err = desc.NewDecoder(tarReader).Decode(&pkg); err != nil {
			return fmt.Errorf("failed to decode %s: %w", hdr.Name, err)
		}

		if pkg.Arch != arch && pkg.Arch != "any" {
			return fmt.Errorf("%s arch mismatch: %s", pkg.Name, pkg.Arch)
		}

		r.addPackage(&pkg)
	}

	return nil
}

func (r *repo) addPackage(p *Package) {
	arch, ok := r.packages[p.Arch]
	if !ok {
		arch = make(map[string]map[string]*Package)
		r.packages[p.Arch] = arch
	}
	pkg, ok := arch[p.Name]
	if !ok {
		pkg = make(map[string]*Package)
		arch[p.Name] = pkg
	}
	pkg[p.Version] = p
}

// Add adds a package to the repository.
func (r *repo) Add(filename string, data []byte) error {
	p, err := parsePkgInfo(filename, bytes.NewReader(data))
	if err != nil {
		return fmt.Errorf("failed to parse .PKGINFO for %s: %w", filename, err)
	}

	p.SHA256Sum = sha256.Sum256(data)
	p.CompressedSize = int64(len(data))
	p.Filename = filename

	pkgPath := filepath.Join(r.dir, p.Arch, filename)
	if err := os.MkdirAll(filepath.Dir(pkgPath), 0o750); err != nil {
		return fmt.Errorf("failed to create arch folder: %w", err)
	}
	if err := os.WriteFile(pkgPath, data, 0o600); err != nil {
		return fmt.Errorf("failed to write package file: %w", err)
	}

	if r.pgpKey != nil {
		sig, err := pgp.Sign(r.pgpKey, data)
		if err != nil {
			return fmt.Errorf("signing package file: %w", err)
		}
		sigPath := pkgPath + ".sig"
		if err := os.WriteFile(sigPath, sig, 0o600); err != nil {
			return fmt.Errorf("failed to write signature: %w", err)
		}
	}

	r.addPackage(p)

	return nil
}

// Write writes the repository to disk. It creates a separate .db file for each arch.
// It also writes the public key to key.asc if it exists.
func (r *repo) Write() error {
	for arch, pkgs := range r.packages {
		if err := r.writeDB(arch, pkgs); err != nil {
			return fmt.Errorf("failed to write db for %s: %w", arch, err)
		}
	}

	if r.pgpKey != nil {
		key, err := pgp.MarshalPublicKey(pgp.Public(r.pgpKey))
		if err != nil {
			return fmt.Errorf("failed to marshal key: %w", err)
		}
		ascPath := filepath.Join(r.dir, "key.asc")
		if err := os.WriteFile(ascPath, key, 0o600); err != nil {
			return fmt.Errorf("failed to write key: %w", err)
		}
	}

	return nil
}

//nolint:funlen
func (r *repo) writeDB(arch string, pkgs map[string]map[string]*Package) error {
	dbPath := filepath.Join(r.dir, arch, r.name+".db")

	if err := os.MkdirAll(filepath.Dir(dbPath), 0o750); err != nil {
		return fmt.Errorf("failed to create folder: %w", err)
	}

	dbFile, err := os.Create(dbPath)
	if err != nil {
		return fmt.Errorf("failed to create db file: %w", err)
	}
	defer dbFile.Close()

	zstdWriter, err := zstd.NewWriter(dbFile,
		zstd.WithEncoderLevel(zstd.SpeedBestCompression),
		zstd.WithZeroFrames(true),
		zstd.WithEncoderCRC(false),
	)
	if err != nil {
		return fmt.Errorf("failed to create zstd writer: %w", err)
	}
	tw := tar.NewWriter(zstdWriter)

	for pkgName, versions := range pkgs {
		var latestPkg *Package
		for _, pkg := range versions {
			if latestPkg == nil || compareVersions(pkg.Version, latestPkg.Version) > 0 {
				latestPkg = pkg
			}
		}

		b, err := desc.Marshal(latestPkg)
		if err != nil {
			return fmt.Errorf("failed to marshal desc for %s: %w", pkgName, err)
		}
		h := &tar.Header{
			Name:    latestPkg.Name + "-" + latestPkg.Version + "/desc",
			Mode:    0o644,
			Size:    int64(len(b)),
			ModTime: time.Unix(0, 0),
			Uid:     0,
			Gid:     0,
		}
		if err := tw.WriteHeader(h); err != nil {
			return fmt.Errorf("failed to write tar header: %w", err)
		}
		if _, err := tw.Write(b); err != nil {
			return fmt.Errorf("failed to write desc: %w", err)
		}
	}

	if err := tw.Close(); err != nil {
		return fmt.Errorf("failed to close tar writer: %w", err)
	}
	if err := zstdWriter.Close(); err != nil {
		return fmt.Errorf("failed to close zstd writer: %w", err)
	}
	if err := dbFile.Close(); err != nil {
		return fmt.Errorf("failed to close db file: %w", err)
	}

	if r.pgpKey != nil {
		dbData, err := os.ReadFile(dbPath)
		if err != nil {
			return fmt.Errorf("failed to read db file: %w", err)
		}
		sig, err := pgp.Sign(r.pgpKey, dbData)
		if err != nil {
			return fmt.Errorf("failed to sign db file: %w", err)
		}
		sigPath := dbPath + ".sig"
		if err := os.WriteFile(sigPath, sig, 0o600); err != nil {
			return fmt.Errorf("failed to write sig: %w", err)
		}
	}

	return nil
}

// parsePkgInfo extracts .PKGINFO from the .pkg.tar.* file.
func parsePkgInfo(filename string, f io.Reader) (*Package, error) {
	r, err := decompress(filepath.Ext(filename))(f)
	if err != nil {
		return nil, err
	}
	defer r.Close()

	tarReader := tar.NewReader(r)

	pkg := &Package{}
	for {
		hdr, err := tarReader.Next()
		if errors.Is(err, io.EOF) {
			break
		}
		if err != nil {
			return nil, fmt.Errorf("failed to read tar: %w", err)
		}
		if hdr.Typeflag != tar.TypeReg || hdr.Name != ".PKGINFO" {
			continue
		}
		infoData, err := io.ReadAll(tarReader)
		if err != nil {
			return nil, fmt.Errorf("failed to read .PKGINFO: %w", err)
		}
		err = pkginfo.Unmarshal(infoData, pkg)
		if err != nil {
			return nil, fmt.Errorf("failed to parse .PKGINFO: %w", err)
		}
	}

	if pkg.Name == "" || pkg.Version == "" || pkg.Arch == "" {
		return nil, errors.New("missing required fields in .PKGINFO")
	}

	return pkg, nil
}
