package apk

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"

	"github.com/abemedia/appcast/pkg/crypto/rsa"
	"github.com/abemedia/appcast/target"
	"gitlab.alpinelinux.org/alpine/go/repository"
)

type repo struct {
	repos map[string]*repository.ApkIndex
	dir   string
}

func openRepo(ctx context.Context, t target.Target) (*repo, error) {
	dir, err := os.MkdirTemp("", "")
	if err != nil {
		return nil, err
	}

	res := &repo{
		repos: make(map[string]*repository.ApkIndex),
		dir:   dir,
	}

	// See https://wiki.alpinelinux.org/wiki/Architecture
	archs := []string{"x86", "x86_64", "armhf", "armv7", "aarch64", "ppc64le", "s390x"}

	for _, arch := range archs {
		r, err := t.NewReader(ctx, arch+"/APKINDEX.tar.gz")
		if err != nil {
			if errors.Is(err, fs.ErrNotExist) {
				continue
			}
			return nil, err
		}
		index, err := repository.IndexFromArchive(r)
		if err != nil {
			return nil, err
		}
		res.repos[arch] = index
	}

	return res, nil
}

func (r *repo) Add(b []byte) error {
	p, err := repository.ParsePackage(bytes.NewReader(b))
	if err != nil {
		return err
	}

	index, ok := r.repos[p.Arch]
	if !ok {
		index = &repository.ApkIndex{}
		r.repos[p.Arch] = index
	}
	index.Packages = append(index.Packages, p)

	dirname := filepath.Join(r.dir, p.Arch)
	if err = os.MkdirAll(dirname, fs.ModePerm); err != nil {
		return err
	}

	filename := fmt.Sprintf("%s-%s.apk", p.Name, p.Version)
	return os.WriteFile(filepath.Join(dirname, filename), b, 0o600)
}

func (r *repo) Write(rsaKey *rsa.PrivateKey, publicKeyName string) error {
	for arch, index := range r.repos {
		rd, err := repository.ArchiveFromIndex(index)
		if err != nil {
			return err
		}

		if rsaKey != nil {
			rd, err = repository.SignArchive(rd, rsaKey, publicKeyName)
			if err != nil {
				return err
			}
		}

		path := filepath.Join(r.dir, arch, "APKINDEX.tar.gz")
		f, err := os.Create(path)
		if err != nil {
			return err
		}
		if _, err = io.Copy(f, rd); err != nil {
			return err
		}
		if err = f.Close(); err != nil {
			return err
		}
	}

	if rsaKey != nil {
		pub, err := rsa.MarshalPublicKey(rsa.Public(rsaKey))
		if err != nil {
			return err
		}
		return os.WriteFile(filepath.Join(r.dir, publicKeyName), pub, 0o600)
	}

	return nil
}
