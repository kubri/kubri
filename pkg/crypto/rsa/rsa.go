package rsa

import (
	stdcrypto "crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha1"
	"crypto/x509"
	"encoding/pem"

	"github.com/kubri/kubri/pkg/crypto"
)

type (
	PrivateKey = rsa.PrivateKey
	PublicKey  = rsa.PublicKey
)

// Sign signs the data with the private key.
func Sign(key *PrivateKey, data []byte) ([]byte, error) {
	if key == nil {
		return nil, crypto.ErrInvalidKey
	}
	hashed := sha1.Sum(data)
	return rsa.SignPKCS1v15(nil, key, stdcrypto.SHA1, hashed[:])
}

// Verify verifies the signature of the data with the public key.
func Verify(key *PublicKey, data, sig []byte) bool {
	if key == nil {
		return false
	}
	hashed := sha1.Sum(data)
	err := rsa.VerifyPKCS1v15(key, stdcrypto.SHA1, hashed[:], sig)
	return err == nil
}

// NewPrivateKey returns a new private key.
func NewPrivateKey() (*PrivateKey, error) {
	return rsa.GenerateKey(rand.Reader, 2048)
}

// MarshalPrivateKey returns the PEM encoded private key.
func MarshalPrivateKey(key *PrivateKey) ([]byte, error) {
	if key == nil {
		return nil, crypto.ErrInvalidKey
	}
	b, err := x509.MarshalPKCS8PrivateKey(key)
	if err != nil {
		return nil, err
	}
	return pem.EncodeToMemory(&pem.Block{Type: "PRIVATE KEY", Bytes: b}), nil
}

// UnmarshalPrivateKey returns a private key from a PEM encoded key.
func UnmarshalPrivateKey(b []byte) (*PrivateKey, error) {
	block, _ := pem.Decode(b)
	if block == nil {
		return nil, crypto.ErrInvalidKey
	}
	key, err := x509.ParsePKCS8PrivateKey(block.Bytes)
	if err != nil {
		return nil, crypto.ErrInvalidKey
	}
	rsaKey, ok := key.(*PrivateKey)
	if !ok {
		return nil, crypto.ErrWrongKeyType
	}
	return rsaKey, nil
}

// Public extracts the public key from a private key.
func Public(key *PrivateKey) *PublicKey {
	return key.Public().(*PublicKey) //nolint:forcetypeassert
}

// MarshalPublicKey returns the PEM encoded public key.
func MarshalPublicKey(key *PublicKey) ([]byte, error) {
	if key == nil {
		return nil, crypto.ErrInvalidKey
	}
	b, err := x509.MarshalPKIXPublicKey(key)
	if err != nil {
		return nil, err
	}
	return pem.EncodeToMemory(&pem.Block{Type: "PUBLIC KEY", Bytes: b}), nil
}

// UnmarshalPublicKey returns a public key from a PEM encoded key.
func UnmarshalPublicKey(b []byte) (*PublicKey, error) {
	block, _ := pem.Decode(b)
	if block == nil {
		return nil, crypto.ErrInvalidKey
	}
	key, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return nil, crypto.ErrInvalidKey
	}
	rsaKey, ok := key.(*PublicKey)
	if !ok {
		return nil, crypto.ErrWrongKeyType
	}
	return rsaKey, nil
}
