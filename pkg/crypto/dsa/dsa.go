package dsa

import (
	"crypto/dsa"
	"crypto/rand"
	"crypto/sha1"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/asn1"
	"encoding/pem"
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

type signature struct {
	R, S *big.Int
}

func Sign(key *PrivateKey, data []byte) ([]byte, error) {
	if key == nil {
		return nil, crypto.ErrInvalidKey
	}
	sum := sha1.Sum(data)
	r, s, err := dsa.Sign(rand.Reader, key, sum[:])
	if err != nil {
		return nil, err
	}
	return asn1.Marshal(signature{r, s})
}

func Verify(key *PublicKey, data, sig []byte) bool {
	if key == nil {
		return false
	}
	var s signature
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
	if key == nil {
		return nil, crypto.ErrInvalidKey
	}
	b, err := asn1.Marshal(privateKey{0, key.P, key.Q, key.G, key.Y, key.X})
	if err != nil {
		return nil, err
	}
	return pem.EncodeToMemory(&pem.Block{Type: "PRIVATE KEY", Bytes: b}), nil
}

func UnmarshalPrivateKey(b []byte) (*PrivateKey, error) {
	block, _ := pem.Decode(b)
	if block == nil {
		return nil, crypto.ErrInvalidKey
	}
	var k privateKey
	if _, err := asn1.Unmarshal(block.Bytes, &k); err != nil {
		return nil, crypto.ErrInvalidKey
	}

	return &dsa.PrivateKey{
		PublicKey: dsa.PublicKey{
			Parameters: dsa.Parameters{P: k.P, Q: k.Q, G: k.G},
			Y:          k.Y,
		},
		X: k.X,
	}, nil
}

func Public(key *PrivateKey) *PublicKey {
	return &key.PublicKey
}

func MarshalPublicKey(key *PublicKey) ([]byte, error) {
	if key == nil {
		return nil, crypto.ErrInvalidKey
	}

	var pub struct {
		Algo      pkix.AlgorithmIdentifier
		BitString asn1.BitString
	}
	pub.Algo.Algorithm = []int{1, 2, 840, 10040, 4, 1}
	pub.Algo.Parameters.FullBytes, _ = asn1.Marshal(key.Parameters)
	pub.BitString.Bytes, _ = asn1.Marshal(key.Y)
	pub.BitString.BitLength = len(pub.BitString.Bytes) * 8

	b, err := asn1.Marshal(pub)
	if err != nil {
		return nil, err
	}

	return pem.EncodeToMemory(&pem.Block{Type: "PUBLIC KEY", Bytes: b}), nil
}

func UnmarshalPublicKey(b []byte) (*PublicKey, error) {
	block, _ := pem.Decode(b)
	if block == nil {
		return nil, crypto.ErrInvalidKey
	}
	key, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return nil, crypto.ErrInvalidKey
	}
	dsaKey, ok := key.(*PublicKey)
	if !ok {
		return nil, crypto.ErrWrongKeyType
	}
	return dsaKey, nil
}
