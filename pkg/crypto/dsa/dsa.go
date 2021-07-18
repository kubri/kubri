package dsa

import (
	"crypto/dsa"
	"crypto/rand"
	"crypto/sha1"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/asn1"
	"fmt"
	"math/big"

	"github.com/abemedia/appcast/pkg/crypto"
)

type (
	PrivateKey = dsa.PrivateKey
	PublicKey  = dsa.PublicKey
)

type privateKey struct {
	Version       int
	P, Q, G, Y, X *big.Int
}

func Sign(key *PrivateKey, data []byte) ([]byte, error) {
	sum := sha1.Sum(data)
	r, s, err := dsa.Sign(rand.Reader, key, sum[:])
	if err != nil {
		return nil, err
	}
	return asn1.Marshal(struct{ R, S *big.Int }{r, s})
}

func Verify(key *PublicKey, data, sig []byte) bool {
	var s struct{ R, S *big.Int }
	if _, err := asn1.Unmarshal(sig, &s); err != nil {
		return false
	}
	sum := sha1.Sum(data)
	return dsa.Verify(key, sum[:], s.R, s.S)
}

func NewPrivateKey() (*PrivateKey, error) {
	var key PrivateKey
	err := dsa.GenerateParameters(&key.Parameters, rand.Reader, dsa.L3072N256)
	if err != nil {
		return nil, err
	}
	err = dsa.GenerateKey(&key, rand.Reader)
	if err != nil {
		return nil, err
	}
	return &key, nil
}

func MarshalPrivateKey(key *PrivateKey) ([]byte, error) {
	return asn1.Marshal(privateKey{
		P: key.P,
		Q: key.Q,
		G: key.G,
		Y: key.Y,
		X: key.X,
	})
}

func UnmarshalPrivateKey(b []byte) (*PrivateKey, error) {
	var k privateKey
	if _, err := asn1.Unmarshal(b, &k); err != nil {
		return nil, fmt.Errorf("failed to parse DSA key: %w", err)
	}

	return &dsa.PrivateKey{
		PublicKey: dsa.PublicKey{
			Parameters: dsa.Parameters{
				P: k.P,
				Q: k.Q,
				G: k.G,
			},
			Y: k.Y,
		},
		X: k.X,
	}, nil
}

func NewPublicKey(key *PrivateKey) *PublicKey {
	return &key.PublicKey
}

func MarshalPublicKey(key *PublicKey) ([]byte, error) {
	var pub struct {
		Algo      pkix.AlgorithmIdentifier
		BitString asn1.BitString
	}
	pub.Algo.Algorithm = []int{1, 2, 840, 10040, 4, 1}
	pub.Algo.Parameters.FullBytes, _ = asn1.Marshal(key.Parameters)
	pub.BitString.Bytes, _ = asn1.Marshal(key.Y)
	pub.BitString.BitLength = len(pub.BitString.Bytes) * 8
	return asn1.Marshal(pub)
}

func UnmarshalPublicKey(b []byte) (*PublicKey, error) {
	key, err := x509.ParsePKIXPublicKey(b)
	if err != nil {
		return nil, err
	}
	dsaKey, ok := key.(*PublicKey)
	if !ok {
		return nil, crypto.ErrWrongKeyType
	}
	return dsaKey, nil
}
