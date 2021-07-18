package ed25519

import (
	"crypto/ed25519"
	"crypto/rand"
	"crypto/x509"

	"github.com/abemedia/appcast/pkg/crypto"
)

type (
	PrivateKey = ed25519.PrivateKey
	PublicKey  = ed25519.PublicKey
)

func Sign(key PrivateKey, data []byte) []byte {
	return ed25519.Sign(key, data)
}

func Verify(key PublicKey, data, sig []byte) bool {
	return ed25519.Verify(key, data, sig)
}

func NewPrivateKey() (PrivateKey, error) {
	_, key, err := ed25519.GenerateKey(rand.Reader)
	return key, err
}

func MarshalPrivateKey(key PrivateKey) ([]byte, error) {
	return x509.MarshalPKCS8PrivateKey(key)
}

func UnmarshalPrivateKey(b []byte) (PrivateKey, error) {
	key, err := x509.ParsePKCS8PrivateKey(b)
	if err != nil {
		return nil, err
	}
	edKey, ok := key.(PrivateKey)
	if !ok {
		return nil, crypto.ErrWrongKeyType
	}
	return edKey, nil
}

func Public(key PrivateKey) PublicKey {
	return key.Public().(PublicKey) //nolint:forcetypeassert
}

func MarshalPublicKey(key PublicKey) ([]byte, error) {
	return x509.MarshalPKIXPublicKey(key)
}

func UnmarshalPublicKey(b []byte) (PublicKey, error) {
	key, err := x509.ParsePKIXPublicKey(b)
	if err != nil {
		return nil, err
	}
	edKey, ok := key.(PublicKey)
	if !ok {
		return nil, crypto.ErrWrongKeyType
	}
	return edKey, nil
}
