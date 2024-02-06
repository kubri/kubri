package ed25519_test

import (
	"crypto/x509"
	"os"
	"os/exec"
	"path/filepath"
	"testing"

	"github.com/kubri/kubri/pkg/crypto/ed25519"
	"github.com/kubri/kubri/pkg/crypto/internal/cryptotest"
)

func TestEd25519(t *testing.T) {
	cryptotest.Test(t, cryptotest.Implementation[ed25519.PrivateKey, ed25519.PublicKey]{
		NewPrivateKey:       ed25519.NewPrivateKey,
		MarshalPrivateKey:   ed25519.MarshalPrivateKey,
		UnmarshalPrivateKey: ed25519.UnmarshalPrivateKey,
		Public:              ed25519.Public,
		MarshalPublicKey:    ed25519.MarshalPublicKey,
		UnmarshalPublicKey:  ed25519.UnmarshalPublicKey,
		Sign:                ed25519.Sign,
		Verify:              ed25519.Verify,
	})

	t.Run("OpenSSL", func(t *testing.T) {
		if _, err := exec.LookPath("openssl"); err != nil {
			t.Skip("openssl not in path")
		}

		priv, _ := ed25519.NewPrivateKey()
		data := []byte("foo\nbar\nbaz")
		sig, _ := ed25519.Sign(priv, data)
		pub, _ := x509.MarshalPKIXPublicKey(ed25519.Public(priv))

		dir := t.TempDir()
		os.WriteFile(filepath.Join(dir, "public.der"), pub, 0o600)
		os.WriteFile(filepath.Join(dir, "data.txt"), data, 0o600)
		os.WriteFile(filepath.Join(dir, "data.txt.sig"), sig, 0o600)

		cmd := exec.Command("openssl", "pkeyutl", "-verify", "-pubin", "-inkey", "public.der", "-rawin", "-in", "data.txt", "-sigfile", "data.txt.sig")
		cmd.Dir = dir
		out, err := cmd.CombinedOutput()
		if err != nil {
			t.Fatal(err, string(out))
		}
		t.Log(string(out))
	})
}
