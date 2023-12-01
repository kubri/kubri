package apk

import (
	"archive/tar"
	"bytes"
	"compress/gzip"
	"crypto"
	"crypto/rsa"
	"crypto/sha1"
	"fmt"
	"io"
)

// TODO: Remove when merged: https://gitlab.alpinelinux.org/alpine/go/-/merge_requests/42
func sign(archive io.Reader, privateKey *rsa.PrivateKey, keyName string) (signedArchive io.Reader, err error) {
	b, err := io.ReadAll(archive)
	if err != nil {
		return nil, err
	}

	var tarballContents bytes.Buffer
	gw := gzip.NewWriter(&tarballContents)
	defer gw.Close()
	tw := tar.NewWriter(gw)

	digest := sha1.Sum(b)
	sig, err := privateKey.Sign(nil, digest[:], crypto.SHA1)
	if err != nil {
		return nil, err
	}

	h := &tar.Header{Name: ".SIGN.RSA." + keyName, Size: int64(len(sig)), Mode: 0o600}
	if err = tw.WriteHeader(h); err != nil {
		return nil, fmt.Errorf("writing tar header for %s: %w", h.Name, err)
	}

	if _, err = tw.Write(sig); err != nil {
		return nil, fmt.Errorf("copying tar contents for %s: %w", h.Name, err)
	}

	if err = tw.Flush(); err != nil {
		return nil, fmt.Errorf("copying tar contents for %s: %w", h.Name, err)
	}

	return io.MultiReader(&tarballContents, bytes.NewReader(b)), nil
}
