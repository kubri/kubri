package arch

import (
	"archive/tar"
	"bytes"
	"context"
	"crypto/sha256"
	"errors"
	"fmt"
	"io"
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

// repo holds the repository state. We track all arches in one place.
// The map is of the form: pkgByArch[arch][[2]string{name,version}] = *Package
// That way, we can easily iterate by arch to produce separate .db files.
type repo struct {
	dir      string
	name     string
	pgpKey   *pgp.PrivateKey
	packages map[string]map[string]map[string]*Package
}

// openRepo creates a temporary directory for staging and iterates over all arch subfolders
// in the target looking for <arch>/<repoName>.db, then parses them into pkgByArch.
// This allows you to parse multiple arches in a single pass if you want.
func openRepo(ctx context.Context, t target.Target, repoName string, pgpKey *pgp.PrivateKey) (*repo, error) {
	dir, err := os.MkdirTemp("", "archrepo-")
	if err != nil {
		return nil, err
	}

	r := &repo{
		dir:      dir,
		name:     repoName,
		pgpKey:   pgpKey,
		packages: make(map[string]map[string]map[string]*Package), // arch -> pkgName -> versions
	}

	// TODO: Implement `Walk` on targets so we can iterate over all folders and
	// avoid hard-coding the arches.
	archs := []string{
		"x86_64",
		"aarch64", "armv7h", // https://archlinuxarm.org/packages
		"powerpc64le", "powerpc64", "powerpc", "riscv64", // https://archlinuxpower.org/
		"i686", "pentium4", // https://archlinux32.org/architecture/
	}

	for _, arch := range archs {
		dbPathInTarget := filepath.Join(arch, repoName+".db")
		existingDB, err := t.NewReader(ctx, dbPathInTarget)
		if err != nil {
			continue
		}
		defer existingDB.Close()

		if err := r.parseExistingDB(existingDB, arch); err != nil {
			return r, fmt.Errorf("failed parsing existing DB %s: %w", dbPathInTarget, err)
		}
	}

	return r, nil
}

func (r *repo) parseExistingDB(dbData io.ReadCloser, arch string) error {
	zr, err := zstd.NewReader(dbData)
	if err != nil {
		return fmt.Errorf("zstd.NewReader: %w", err)
	}
	defer zr.Close()

	tarReader := tar.NewReader(zr)

	for {
		hdr, err := tarReader.Next()
		if errors.Is(err, io.EOF) {
			break
		}
		if err != nil {
			return fmt.Errorf("reading tar: %w", err)
		}
		if hdr.Typeflag != tar.TypeReg {
			continue
		}
		if !strings.HasSuffix(hdr.Name, "/desc") {
			continue
		}

		var pkg Package
		if err = desc.NewDecoder(tarReader).Decode(&pkg); err != nil {
			return fmt.Errorf("decoding desc: %w", err)
		}

		// If the arch from desc doesn't match, skip
		if pkg.Arch != arch {
			continue
		}

		r.storePackage(&pkg)
	}

	return nil
}

func (r *repo) storePackage(pkg *Package) {
	archMap, ok := r.packages[pkg.Arch]
	if !ok {
		archMap = make(map[string]map[string]*Package)
		r.packages[pkg.Arch] = archMap
	}
	pkgMap, ok := archMap[pkg.Name]
	if !ok {
		pkgMap = make(map[string]*Package)
		archMap[pkg.Name] = pkgMap
	}
	pkgMap[pkg.Version] = pkg
}

// Add adds a package to the repository.
func (r *repo) Add(filename string, data []byte) error {
	p, err := parsePkgInfo(filename, bytes.NewReader(data))
	if err != nil {
		return fmt.Errorf("failed to parse .PKGINFO: %w", err)
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
			return fmt.Errorf("writing sig: %w", err)
		}
	}

	r.storePackage(p)
	return nil
}

// Write writes the repository to disk. It creates a separate .db file for each arch.
// It also writes the public key to key.asc if it exists.
func (r *repo) Write() error {
	for arch, pkgMap := range r.packages {
		if err := r.writeDB(arch, pkgMap); err != nil {
			return fmt.Errorf("writing arch DB for %s: %w", arch, err)
		}
	}

	if r.pgpKey != nil {
		key, err := pgp.MarshalPublicKey(pgp.Public(r.pgpKey))
		if err != nil {
			return fmt.Errorf("exporting public key: %w", err)
		}
		ascPath := filepath.Join(r.dir, "key.asc")
		if err := os.WriteFile(ascPath, key, 0o600); err != nil {
			return fmt.Errorf("writing public key: %w", err)
		}
	}

	return nil
}

//nolint:funlen
func (r *repo) writeDB(arch string, pkgMap map[string]map[string]*Package) error {
	dbPath := filepath.Join(r.dir, arch, r.name+".db")

	if err := os.MkdirAll(filepath.Dir(dbPath), 0o750); err != nil {
		return fmt.Errorf("creating arch folder: %w", err)
	}

	dbFile, err := os.Create(dbPath)
	if err != nil {
		return fmt.Errorf("creating db file: %w", err)
	}
	defer dbFile.Close()

	zstdWriter, err := zstd.NewWriter(dbFile,
		zstd.WithEncoderLevel(zstd.SpeedBestCompression),
		zstd.WithZeroFrames(true),
		zstd.WithEncoderCRC(false),
	)
	if err != nil {
		return fmt.Errorf("creating zstd writer: %w", err)
	}
	tw := tar.NewWriter(zstdWriter)

	for pkgName, versions := range pkgMap {
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
		descHeader := &tar.Header{
			Name:    latestPkg.Name + "-" + latestPkg.Version + "/desc",
			Mode:    0o644,
			Size:    int64(len(b)),
			ModTime: time.Unix(0, 0),
			Uid:     0,
			Gid:     0,
		}
		if err := tw.WriteHeader(descHeader); err != nil {
			return fmt.Errorf("failed to write tar header: %w", err)
		}
		if _, err := tw.Write(b); err != nil {
			return fmt.Errorf("failed to write desc: %w", err)
		}
	}

	// Close tar and zstd writers.
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
			return fmt.Errorf("reading db file for signing: %w", err)
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
	var tarReader *tar.Reader
	switch {
	case strings.HasSuffix(filename, ".pkg.tar.zst"):
		zr, err := zstd.NewReader(f)
		if err != nil {
			return nil, fmt.Errorf("zstd reader: %w", err)
		}
		defer zr.Close()
		tarReader = tar.NewReader(zr)
	case strings.HasSuffix(filename, ".pkg.tar.gz"):
		return nil, errors.New("TODO: implement .pkg.tar.gz reading")
	case strings.HasSuffix(filename, ".pkg.tar.xz"):
		return nil, errors.New("TODO: implement .pkg.tar.xz reading")
	default:
		return nil, fmt.Errorf("unhandled package compression: %s", filename)
	}

	pkg := &Package{}
	for {
		hdr, err := tarReader.Next()
		if errors.Is(err, io.EOF) {
			break
		}
		if err != nil {
			return nil, fmt.Errorf("failed to read tar: %w", err)
		}
		if hdr.Typeflag != tar.TypeReg {
			continue
		}
		if hdr.Name == ".PKGINFO" {
			infoData, err := io.ReadAll(tarReader)
			if err != nil {
				return nil, fmt.Errorf("failed to read .PKGINFO: %w", err)
			}
			err = pkginfo.Unmarshal(infoData, pkg)
			if err != nil {
				return nil, fmt.Errorf("failed to parse .PKGINFO: %w", err)
			}
			break
		}
	}

	if pkg.Name == "" || pkg.Version == "" || pkg.Arch == "" {
		return nil, errors.New("missing required fields in .PKGINFO")
	}
	return pkg, nil
}
