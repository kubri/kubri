package ed25519_test

import (
	"bytes"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"testing"

	"github.com/kubri/kubri/pkg/crypto"
	"github.com/kubri/kubri/pkg/crypto/ed25519"
)

func TestPEM(t *testing.T) {
	priv, _ := ed25519.NewPrivateKey()
	var keyPEM []byte

	t.Run("MarshalPrivateKeyPEM", func(t *testing.T) {
		b, err := ed25519.MarshalPrivateKeyPEM(priv)
		if err != nil {
			t.Fatal(err)
		}
		keyPEM = b

		if _, err = ed25519.MarshalPrivateKeyPEM(nil); err != crypto.ErrInvalidKey {
			t.Fatalf("expected %v, got %v", crypto.ErrInvalidKey, err)
		}

		if _, err = ed25519.MarshalPrivateKeyPEM(priv[:32]); err != crypto.ErrInvalidKey {
			t.Fatalf("expected %v, got %v", crypto.ErrInvalidKey, err)
		}
	})

	t.Run("UnmarshalPrivateKeyPEM", func(t *testing.T) {
		key, err := ed25519.UnmarshalPrivateKeyPEM(keyPEM)
		if err != nil {
			t.Fatal(err)
		}
		if !bytes.Equal(key, priv) {
			t.Fatal("keys do not match")
		}

		if _, err = ed25519.UnmarshalPrivateKeyPEM(nil); err != crypto.ErrInvalidKey {
			t.Fatalf("expected %v, got %v", crypto.ErrInvalidKey, err)
		}

		if _, err = ed25519.UnmarshalPrivateKeyPEM([]byte("invalid")); err != crypto.ErrInvalidKey {
			t.Fatalf("expected %v, got %v", crypto.ErrInvalidKey, err)
		}

		invalidPEM := pem.EncodeToMemory(&pem.Block{Type: "PRIVATE KEY", Bytes: []byte("invalid")})
		if _, err = ed25519.UnmarshalPrivateKeyPEM(invalidPEM); err != crypto.ErrInvalidKey {
			t.Fatalf("expected %v, got %v", crypto.ErrInvalidKey, err)
		}

		rsaPriv, _ := rsa.GenerateKey(rand.Reader, 2048)
		b, _ := x509.MarshalPKCS8PrivateKey(rsaPriv)
		invalidPEM = pem.EncodeToMemory(&pem.Block{Type: "PRIVATE KEY", Bytes: b})
		if _, err = ed25519.UnmarshalPrivateKeyPEM(invalidPEM); err != crypto.ErrWrongKeyType {
			t.Fatalf("expected %v, got %v", crypto.ErrInvalidKey, err)
		}
	})
}
