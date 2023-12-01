package rsa

import (
	stdcrypto "crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha1"
	"crypto/x509"
	"encoding/pem"

	"github.com/abemedia/appcast/pkg/crypto"
)

type (
	PrivateKey = rsa.PrivateKey
	PublicKey  = rsa.PublicKey
)

func Sign(key *PrivateKey, data []byte) ([]byte, error) {
	if key == nil {
		return nil, crypto.ErrInvalidKey
	}
	hashed := sha1.Sum(data)
	return rsa.SignPKCS1v15(nil, key, stdcrypto.SHA1, hashed[:])
}

func Verify(key *PublicKey, data, sig []byte) bool {
	if key == nil {
		return false
	}
	hashed := sha1.Sum(data)
	err := rsa.VerifyPKCS1v15(key, stdcrypto.SHA1, hashed[:], sig)
	return err == nil
}

func NewPrivateKey() (*PrivateKey, error) {
	return rsa.GenerateKey(rand.Reader, 2048)
}

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

func Public(key *PrivateKey) *PublicKey {
	return key.Public().(*PublicKey) //nolint:forcetypeassert
}

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
