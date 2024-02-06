package dsa_test

import (
	"math/big"
	"testing"

	"github.com/google/go-cmp/cmp"

	"github.com/kubri/kubri/pkg/crypto/dsa"
	"github.com/kubri/kubri/pkg/crypto/internal/cryptotest"
)

func TestDSA(t *testing.T) {
	cryptotest.Test(t,
		cryptotest.Implementation[*dsa.PrivateKey, *dsa.PublicKey]{
			NewPrivateKey:       dsa.NewPrivateKey,
			MarshalPrivateKey:   dsa.MarshalPrivateKey,
			UnmarshalPrivateKey: dsa.UnmarshalPrivateKey,
			Public:              dsa.Public,
			MarshalPublicKey:    dsa.MarshalPublicKey,
			UnmarshalPublicKey:  dsa.UnmarshalPublicKey,
			Sign:                dsa.Sign,
			Verify:              dsa.Verify,
		},
		cryptotest.WithCmpOptions(cmp.Comparer(func(a, b *big.Int) bool { return a.Cmp(b) == 0 })),
		cryptotest.WithOpenSSLTest("dgst", "-sha1", "-verify", "public.pem", "-signature", "data.txt.sig", "data.txt"),
	)
}
