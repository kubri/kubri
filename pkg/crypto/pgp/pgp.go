package pgp

import (
	"bytes"
	"errors"
	"fmt"

	pgperrors "github.com/ProtonMail/go-crypto/openpgp/errors"
	"github.com/ProtonMail/gopenpgp/v2/armor"
	"github.com/ProtonMail/gopenpgp/v2/constants"
	pgpcrypto "github.com/ProtonMail/gopenpgp/v2/crypto"

	"github.com/kubri/kubri/pkg/crypto"
)

type (
	PrivateKey = pgpcrypto.Key
	PublicKey  = pgpcrypto.Key
)

// NewPrivateKey returns a new private key.
func NewPrivateKey(name, email string) (*PrivateKey, error) {
	return pgpcrypto.GenerateKey(name, email, "x25519", 0)
}

// MarshalPrivateKey returns the armored private key.
func MarshalPrivateKey(key *PrivateKey) ([]byte, error) {
	if key == nil {
		return nil, crypto.ErrInvalidKey
	}
	if !key.IsPrivate() {
		return nil, crypto.ErrWrongKeyType
	}
	s, err := key.ArmorWithCustomHeaders("", "")
	if err != nil {
		return nil, err
	}
	return []byte(s), nil
}

// UnmarshalPrivateKey returns a private key from an armored key.
func UnmarshalPrivateKey(b []byte) (*PrivateKey, error) {
	key, err := pgpcrypto.NewKeyFromArmoredReader(bytes.NewReader(b))
	if err != nil {
		return nil, wrapError(crypto.ErrInvalidKey, err)
	}
	if !key.IsPrivate() {
		return nil, fmt.Errorf("%w: public key supplied instead of private key", crypto.ErrInvalidKey)
	}
	return key, nil
}

// Public extracts the public key from a private key.
func Public(key *PrivateKey) *PublicKey {
	pub, err := key.ToPublic()
	if err != nil {
		return key
	}
	return pub
}

// MarshalPublicKey returns the armored public key.
func MarshalPublicKey(key *PublicKey) ([]byte, error) {
	if key == nil {
		return nil, crypto.ErrInvalidKey
	}
	if key.IsPrivate() {
		return nil, crypto.ErrWrongKeyType
	}
	s, err := key.GetArmoredPublicKeyWithCustomHeaders("", "")
	if err != nil {
		return nil, err
	}
	return []byte(s), nil
}

// UnmarshalPublicKey returns a public key from an armored key.
func UnmarshalPublicKey(b []byte) (*PublicKey, error) {
	key, err := pgpcrypto.NewKeyFromArmoredReader(bytes.NewReader(b))
	if err != nil {
		return nil, wrapError(crypto.ErrInvalidKey, err)
	}
	if key.IsPrivate() {
		return nil, fmt.Errorf("%w: private key supplied instead of public key", crypto.ErrInvalidKey)
	}
	return key, nil
}

// Sign signs the data with the private key.
func Sign(key *PrivateKey, data []byte) ([]byte, error) {
	return sign(key, data, false)
}

// SignText signs the data with the private key and wraps it in a signed message.
// Data is considered text and canonicalised with CRLF line endings.
func SignText(key *PrivateKey, data []byte) ([]byte, error) {
	data = bytes.ReplaceAll(data, lf, crlf)
	sig, err := sign(key, data, true)
	if err != nil {
		return nil, err
	}
	b := make([]byte, 0, len(startText)+len(data)+len("\r\n")+len(sig))
	b = append(b, startText...)
	b = append(b, data...)
	b = append(b, "\r\n"...)
	b = append(b, sig...)
	return b, nil
}

func sign(key *PrivateKey, data []byte, text bool) ([]byte, error) {
	if key == nil {
		return nil, crypto.ErrInvalidKey
	}
	if !key.IsPrivate() {
		return nil, crypto.ErrWrongKeyType
	}

	// TODO: Unlock locked key using env var passphrase.

	keyring, err := pgpcrypto.NewKeyRing(key)
	if err != nil {
		return nil, err
	}

	msg := pgpcrypto.NewPlainMessage(data)
	msg.TextType = text

	signature, err := keyring.SignDetached(msg)
	if err != nil {
		return nil, err
	}

	sig, err := armor.ArmorWithTypeAndCustomHeaders(signature.Data, constants.PGPSignatureHeader, "", "")
	if err != nil {
		return nil, err
	}

	return []byte(sig), nil
}

// Verify verifies the signature of the data with the public key.
func Verify(key *PublicKey, data, sig []byte) bool {
	signature, err := pgpcrypto.NewPGPSignatureFromArmored(string(sig))
	if err != nil {
		return false
	}

	// TODO: Unlock locked key using env var passphrase.

	keyring, err := pgpcrypto.NewKeyRing(key)
	if err != nil {
		return false
	}

	msg := pgpcrypto.NewPlainMessage(data)

	err = keyring.VerifyDetached(msg, signature, pgpcrypto.GetUnixTime())
	return err == nil
}

// Split splits a signed message into data and signature.
func Split(msg []byte) (data, sig []byte, _ error) {
	start := bytes.Index(msg, startText)
	end := bytes.Index(msg, endText)

	if start == -1 || end == -1 {
		return nil, nil, ErrInvalidMessage
	}

	return bytes.ReplaceAll(msg[start+len(startText):end], crlf, lf), msg[end+2:], nil
}

//nolint:gochecknoglobals
var (
	startText = []byte("-----BEGIN PGP SIGNED MESSAGE-----\r\nHash: SHA512\r\n\r\n")
	endText   = []byte("\r\n-----BEGIN PGP SIGNATURE-----")
	lf        = []byte("\n")
	crlf      = []byte("\r\n")
)

var ErrInvalidMessage = errors.New("pgp: invalid message")

func wrapError(wrapErr, err error) error {
	var e pgperrors.InvalidArgumentError
	if errors.As(err, &e) {
		return fmt.Errorf("%w: %s", wrapErr, string(e))
	}
	return err
}
