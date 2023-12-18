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
	return marshal(key, ed25519.PrivateKeySize)
}

func UnmarshalPrivateKey(b []byte) (PrivateKey, error) {
	return unmarshal(b, ed25519.PrivateKeySize)
}

func Public(key PrivateKey) PublicKey {
	return key.Public().(PublicKey) //nolint:forcetypeassert
}

func MarshalPublicKey(key PublicKey) ([]byte, error) {
	return marshal(key, ed25519.PublicKeySize)
}

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
