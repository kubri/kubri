package ed25519

import (
	"crypto/ed25519"
	"crypto/rand"
	"crypto/x509"
	"encoding/pem"

	"github.com/abemedia/appcast/pkg/crypto"
)

type (
	PrivateKey = ed25519.PrivateKey
	PublicKey  = ed25519.PublicKey
)

func Sign(key PrivateKey, data []byte) ([]byte, error) {
	if l := len(key); l != ed25519.PrivateKeySize {
		return nil, crypto.ErrInvalidKey
	}
	return ed25519.Sign(key, data), nil
}

func Verify(key PublicKey, data, sig []byte) bool {
	if l := len(key); l != ed25519.PublicKeySize {
		return false
	}
	return ed25519.Verify(key, data, sig)
}

func NewPrivateKey() (PrivateKey, error) {
	_, key, err := ed25519.GenerateKey(rand.Reader)
	return key, err
}

func MarshalPrivateKey(key PrivateKey) ([]byte, error) {
	if l := len(key); l != ed25519.PrivateKeySize {
		return nil, crypto.ErrInvalidKey
	}
	b, err := x509.MarshalPKCS8PrivateKey(key)
	if err != nil {
		return nil, err
	}
	return pem.EncodeToMemory(&pem.Block{Type: "PRIVATE KEY", Bytes: b}), nil
}

func UnmarshalPrivateKey(b []byte) (PrivateKey, error) {
	block, _ := pem.Decode(b)
	if block == nil {
		return nil, crypto.ErrInvalidKey
	}
	key, err := x509.ParsePKCS8PrivateKey(block.Bytes)
	if err != nil {
		return nil, crypto.ErrInvalidKey
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
	if l := len(key); l != ed25519.PublicKeySize {
		return nil, crypto.ErrInvalidKey
	}
	b, err := x509.MarshalPKIXPublicKey(key)
	if err != nil {
		return nil, err
	}
	return pem.EncodeToMemory(&pem.Block{Type: "PUBLIC KEY", Bytes: b}), nil
}

func UnmarshalPublicKey(b []byte) (PublicKey, error) {
	block, _ := pem.Decode(b)
	if block == nil {
		return nil, crypto.ErrInvalidKey
	}
	key, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return nil, crypto.ErrInvalidKey
	}
	edKey, ok := key.(PublicKey)
	if !ok {
		return nil, crypto.ErrWrongKeyType
	}
	return edKey, nil
}
