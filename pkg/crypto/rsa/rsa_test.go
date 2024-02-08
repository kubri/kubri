package rsa_test

import (
	"testing"

	"github.com/google/go-cmp/cmp"

	"github.com/kubri/kubri/internal/test"
	"github.com/kubri/kubri/pkg/crypto/internal/cryptotest"
	"github.com/kubri/kubri/pkg/crypto/rsa"
)

func TestRSA(t *testing.T) {
	cryptotest.Test(t,
		cryptotest.Implementation[*rsa.PrivateKey, *rsa.PublicKey]{
			NewPrivateKey:       rsa.NewPrivateKey,
			MarshalPrivateKey:   rsa.MarshalPrivateKey,
			UnmarshalPrivateKey: rsa.UnmarshalPrivateKey,
			Public:              rsa.Public,
			MarshalPublicKey:    rsa.MarshalPublicKey,
			UnmarshalPublicKey:  rsa.UnmarshalPublicKey,
			Sign:                rsa.Sign,
			Verify:              rsa.Verify,
		},
		cryptotest.WithCmpOptions(test.CompareRSAPrivateKeys()),
		cryptotest.WithCmpOptions(cmp.Comparer(func(a, b *rsa.PublicKey) bool {
			if a == nil || b == nil {
				return a == b
			}
			return a.Equal(b)
		})),
		cryptotest.WithOpenSSLTest("dgst", "-sha1", "-verify", "public.pem", "-signature", "data.txt.sig", "data.txt"))
}
