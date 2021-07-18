package dsa_test

import (
	"testing"

	"github.com/abemedia/appcast/pkg/crypto/dsa"
)

func TestDSA(t *testing.T) {
	priv, err := dsa.NewPrivateKey()
	if err != nil {
		t.Fatal(err)
	}

	privBytes, err := dsa.MarshalPrivateKey(priv)
	if err != nil {
		t.Fatal(err)
	}

	pub := dsa.Public(priv)
	pubBytes, err := dsa.MarshalPublicKey(pub)
	if err != nil {
		t.Fatal(err)
	}

	priv, err = dsa.UnmarshalPrivateKey(privBytes)
	if err != nil {
		t.Fatal(err)
	}

	pub, err = dsa.UnmarshalPublicKey(pubBytes)
	if err != nil {
		t.Fatal(err)
	}

	data := []byte("test")

	sig, err := dsa.Sign(priv, data)
	if err != nil {
		t.Fatal(err)
	}
	if !dsa.Verify(pub, data, sig) {
		t.Fatal("invalid signature")
	}
}
