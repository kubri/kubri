package appcast_test

import (
	"encoding/base64"
	"encoding/pem"
	"os"
	"os/exec"
	"testing"

	"github.com/abemedia/appcast"
	"github.com/abemedia/appcast/pkg/crypto/dsa"
	"github.com/abemedia/appcast/pkg/crypto/ed25519"
	"github.com/abemedia/appcast/source"
	"github.com/abemedia/appcast/source/memory"
)

func TestSign(t *testing.T) {
	data := []byte("test")

	s, _ := memory.New(source.Config{})
	s.UploadAsset("v1.0.0", "test.dmg", data)
	s.UploadAsset("v1.0.0", "test_64-bit.msi", data)

	dsaKey, _ := dsa.NewPrivateKey()
	b, _ := dsa.MarshalPrivateKey(dsaKey)
	dsaPEM := pem.EncodeToMemory(&pem.Block{Type: "PUBLIC KEY", Bytes: b})

	edKey, _ := ed25519.NewPrivateKey()
	b, _ = ed25519.MarshalPrivateKey(edKey)
	edPEM := pem.EncodeToMemory(&pem.Block{Type: "PUBLIC KEY", Bytes: b})

	c := &appcast.Config{
		Source: s,
		DSAKey: saveData(t, dsaPEM),
		EdKey:  saveData(t, edPEM),
	}

	if err := appcast.Sign(c); err != nil {
		t.Fatal(err)
	}

	b, _ = s.DownloadAsset("v1.0.0", "signatures.txt")
	sigs := appcast.Signatures{}
	sigs.UnmarshalText(b)

	edSig, _ := base64.RawStdEncoding.DecodeString(sigs.Get("test.dmg", "ed25519"))
	testEd25519(t, edKey, data, edSig)

	dsaSig, _ := base64.RawStdEncoding.DecodeString(sigs.Get("test_64-bit.msi", "dsa"))
	testDSA(t, dsaKey, data, dsaSig)
}

func testDSA(t *testing.T, key *dsa.PrivateKey, data, sig []byte) {
	t.Helper()

	pub := dsa.NewPublicKey(key)
	if !dsa.Verify(pub, data, sig) {
		t.Error("invalid signature")
	}

	t.Run("openssl-dgst", func(t *testing.T) {
		if _, err := exec.LookPath("openssl"); err != nil {
			t.Skip("openssl not available")
			return
		}

		b, _ := dsa.MarshalPublicKey(pub)
		pem := pem.EncodeToMemory(&pem.Block{Type: "PUBLIC KEY", Bytes: b})

		cmd := exec.Command("openssl", "dgst", "-verify", saveData(t, pem), "-keyform",
			"PEM", "-sha1", "-signature", saveData(t, sig), "-binary", saveData(t, data))
		if b, err := cmd.CombinedOutput(); err != nil {
			t.Fatalf("%s: %s", err, b)
		}
	})
}

func testEd25519(t *testing.T, key ed25519.PrivateKey, data, sig []byte) {
	t.Helper()

	edPub := ed25519.NewPublicKey(key)
	if !ed25519.Verify(edPub, data, sig) {
		t.Error("invalid signature")
	}

	t.Run("openssl-pkeyutl", func(t *testing.T) {
		if _, err := exec.LookPath("openssl"); err != nil {
			t.Skip("openssl not available")
			return
		}

		b, _ := ed25519.MarshalPublicKey(edPub)
		pem := pem.EncodeToMemory(&pem.Block{Type: "PUBLIC KEY", Bytes: b})

		cmd := exec.Command("openssl", "pkeyutl", "-verify", "-pubin", "-inkey", saveData(t, pem),
			"-rawin", "-in", saveData(t, data), "-sigfile", saveData(t, sig))
		if b, err := cmd.CombinedOutput(); err != nil {
			t.Fatalf("%s: %s", err, b)
		}
	})
}

func saveData(t *testing.T, b []byte) string {
	t.Helper()

	f, err := os.CreateTemp(t.TempDir(), "")
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()

	if _, err = f.Write(b); err != nil {
		t.Fatal(err)
	}

	return f.Name()
}
