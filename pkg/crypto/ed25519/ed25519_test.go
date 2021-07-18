package ed25519_test

import (
	"testing"

	"github.com/abemedia/appcast/pkg/crypto/ed25519"
)

func TestEd25519(t *testing.T) {
	priv, err := ed25519.NewPrivateKey()
	if err != nil {
		t.Fatal(err)
	}

	privBytes, err := ed25519.MarshalPrivateKey(priv)
	if err != nil {
		t.Fatal(err)
	}

	pub := ed25519.Public(priv)
	pubBytes, err := ed25519.MarshalPublicKey(pub)
	if err != nil {
		t.Fatal(err)
	}

	priv, err = ed25519.UnmarshalPrivateKey(privBytes)
	if err != nil {
		t.Fatal(err)
	}

	pub, err = ed25519.UnmarshalPublicKey(pubBytes)
	if err != nil {
		t.Fatal(err)
	}

	data := []byte("test")

	sig := ed25519.Sign(priv, data)
	if !ed25519.Verify(pub, data, sig) {
		t.Fatal("invalid signature")
	}
}
