package ed25519

import (
	"crypto/ed25519"
	"crypto/rand"
	"encoding/base64"

	"github.com/abemedia/appcast/pkg/crypto"
)

type (
	PrivateKey = ed25519.PrivateKey
	PublicKey  = ed25519.PublicKey
)

// Sign signs the data with the private key.
func Sign(key PrivateKey, data []byte) ([]byte, error) {
	if l := len(key); l != ed25519.PrivateKeySize {
		return nil, crypto.ErrInvalidKey
	}
	return ed25519.Sign(key, data), nil
}

// Verify verifies the signature of the data with the public key.
func Verify(key PublicKey, data, sig []byte) bool {
	if l := len(key); l != ed25519.PublicKeySize {
		return false
	}
	return ed25519.Verify(key, data, sig)
}

// NewPrivateKey returns a new private key.
func NewPrivateKey() (PrivateKey, error) {
	_, key, err := ed25519.GenerateKey(rand.Reader)
	return key, err
}

// MarshalPrivateKey returns the base64 encoded private key.
func MarshalPrivateKey(key PrivateKey) ([]byte, error) {
	return marshal(key, ed25519.PrivateKeySize)
}

// UnmarshalPrivateKey returns a private key from a base64 encoded key.
func UnmarshalPrivateKey(b []byte) (PrivateKey, error) {
	return unmarshal(b, ed25519.PrivateKeySize)
}

// Public extracts the public key from a private key.
func Public(key PrivateKey) PublicKey {
	return key.Public().(PublicKey) //nolint:forcetypeassert
}

// MarshalPublicKey returns the base64 encoded public key.
func MarshalPublicKey(key PublicKey) ([]byte, error) {
	return marshal(key, ed25519.PublicKeySize)
}

// UnmarshalPublicKey returns a public key from a base64 encoded key.
func UnmarshalPublicKey(b []byte) (PublicKey, error) {
	return unmarshal(b, ed25519.PublicKeySize)
}

func marshal(key []byte, size int) ([]byte, error) {
	if len(key) != size {
		return nil, crypto.ErrInvalidKey
	}
	b := make([]byte, base64.StdEncoding.EncodedLen(size))
	base64.StdEncoding.Encode(b, key)
	return b, nil
}

func unmarshal(b []byte, size int) ([]byte, error) {
	key := make([]byte, base64.StdEncoding.DecodedLen(len(b)))
	n, err := base64.StdEncoding.Decode(key, b)
	if err != nil || n != size {
		return nil, crypto.ErrInvalidKey
	}
	return key[:n], nil
}
