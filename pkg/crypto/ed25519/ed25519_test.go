package ed25519_test

import (
	"testing"

	"github.com/abemedia/appcast/pkg/crypto/ed25519"
	"github.com/abemedia/appcast/pkg/crypto/internal/cryptotest"
)

func TestEd25519(t *testing.T) {
	cryptotest.Test(t,
		cryptotest.Implementation[ed25519.PrivateKey, ed25519.PublicKey]{
			NewPrivateKey:       ed25519.NewPrivateKey,
			MarshalPrivateKey:   ed25519.MarshalPrivateKey,
			UnmarshalPrivateKey: ed25519.UnmarshalPrivateKey,
			Public:              ed25519.Public,
			MarshalPublicKey:    ed25519.MarshalPublicKey,
			UnmarshalPublicKey:  ed25519.UnmarshalPublicKey,
			Sign:                ed25519.Sign,
			Verify:              ed25519.Verify,
		},
		cryptotest.WithOpenSSLTest("pkeyutl", "-verify", "-pubin", "-inkey", "public.pem", "-rawin", "-in", "data.txt", "-sigfile", "data.txt.sig"))
}
