package pgp

import (
	"bytes"
	"errors"
	"fmt"

	pgperrors "github.com/ProtonMail/go-crypto/openpgp/errors"
	"github.com/ProtonMail/gopenpgp/v2/armor"
	"github.com/ProtonMail/gopenpgp/v2/constants"
	pgpcrypto "github.com/ProtonMail/gopenpgp/v2/crypto"
	"github.com/abemedia/appcast/pkg/crypto"
)

type (
	PrivateKey = pgpcrypto.Key
	PublicKey  = pgpcrypto.Key
)

func NewPrivateKey(name, email string) (*PrivateKey, error) {
	return pgpcrypto.GenerateKey(name, email, "x25519", 0)
}

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

func Public(key *PrivateKey) *PublicKey {
	pub, err := key.ToPublic()
	if err != nil {
		return key
	}
	return pub
}

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

func Sign(key *PrivateKey, data []byte) ([]byte, error) {
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
	msg.TextType = true

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
	msg.TextType = true

	err = keyring.VerifyDetached(msg, signature, pgpcrypto.GetUnixTime())
	return err == nil
}

func Join(data, sig []byte) []byte {
	b := make([]byte, 0, len(startText)+len(data)+len("\n")+len(sig))
	b = append(b, startText...)
	b = append(b, data...)
	b = append(b, "\n"...)
	b = append(b, sig...)
	return b
}

func Split(msg []byte) (data, sig []byte, _ error) {
	start := bytes.Index(msg, startText)
	end := bytes.Index(msg, endText)

	if start == -1 || end == -1 {
		return nil, nil, ErrInvalidMessage
	}

	return msg[start+len(startText) : end], msg[end+1:], nil
}

//nolint:gochecknoglobals
var (
	startText = []byte("-----BEGIN PGP SIGNED MESSAGE-----\nHash: SHA512\n\n")
	endText   = []byte("\n-----BEGIN PGP SIGNATURE-----")
)

var ErrInvalidMessage = errors.New("pgp: invalid message")

func wrapError(wrapErr, err error) error {
	var e pgperrors.InvalidArgumentError
	if errors.As(err, &e) {
		return fmt.Errorf("%w: %s", wrapErr, e)
	}
	return err
}