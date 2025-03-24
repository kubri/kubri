package ed25519

import (
	"crypto/ed25519"
	"crypto/x509"
	"encoding/pem"

	"github.com/kubri/kubri/pkg/crypto"
)

// MarshalPrivateKeyPEM returns the PEM encoded private key.
func MarshalPrivateKeyPEM(key PrivateKey) ([]byte, error) {
	if len(key) != ed25519.PrivateKeySize {
		return nil, crypto.ErrInvalidKey
	}
	b, err := x509.MarshalPKCS8PrivateKey(key)
	if err != nil {
		return nil, crypto.ErrInvalidKey
	}
	return pem.EncodeToMemory(&pem.Block{Type: "PRIVATE KEY", Bytes: b}), nil
}

// UnmarshalPrivateKeyPEM returns a private key from a PEM encoded key.
func UnmarshalPrivateKeyPEM(b []byte) (PrivateKey, error) {
	block, _ := pem.Decode(b)
	if block == nil {
		return nil, crypto.ErrInvalidKey
	}
	key, err := x509.ParsePKCS8PrivateKey(block.Bytes)
	if err != nil {
		return nil, crypto.ErrInvalidKey
	}
	edKey, ok := key.(PrivateKey)
	if !ok || len(edKey) != ed25519.PrivateKeySize {
		return nil, crypto.ErrWrongKeyType
	}
	return edKey, nil
}
