package pgp_test

import (
	"errors"
	"os"
	"os/exec"
	"path/filepath"
	"testing"

	"github.com/google/go-cmp/cmp"

	"github.com/kubri/kubri/internal/test"
	"github.com/kubri/kubri/pkg/crypto"
	"github.com/kubri/kubri/pkg/crypto/internal/cryptotest"
	"github.com/kubri/kubri/pkg/crypto/pgp"
)

func TestPGP(t *testing.T) {
	cryptotest.Test(t,
		cryptotest.Implementation[*pgp.PrivateKey, *pgp.PublicKey]{
			NewPrivateKey: func() (*pgp.PrivateKey, error) {
				return pgp.NewPrivateKey("test", "test@example.com")
			},
			MarshalPrivateKey:   pgp.MarshalPrivateKey,
			UnmarshalPrivateKey: pgp.UnmarshalPrivateKey,
			Public:              pgp.Public,
			MarshalPublicKey:    pgp.MarshalPublicKey,
			UnmarshalPublicKey:  pgp.UnmarshalPublicKey,
			Sign:                pgp.Sign,
			Verify:              pgp.Verify,
		},
		cryptotest.WithCmpOptions(test.ComparePGPKeys()),
	)

	priv, _ := pgp.NewPrivateKey("test", "test@example.com")
	pub := pgp.Public(priv)
	pubBytes, _ := pgp.MarshalPublicKey(pub)
	data := []byte("foo\nbar\nbaz")
	sig, _ := pgp.Sign(priv, data)
	signed, _ := pgp.SignText(priv, data)

	t.Run("NewPrivateKey", func(t *testing.T) {
		tests := []struct {
			desc  string
			name  string
			email string
			err   bool
		}{
			{
				desc: "name only",
				name: "test",
			},
			{
				desc:  "email only",
				email: "test",
			},
			{
				desc: "missing name & email",
				err:  true,
			},
		}

		for _, test := range tests {
			_, err := pgp.NewPrivateKey(test.name, test.email)
			if (err == nil) == test.err {
				t.Errorf("%s should return error %t got %t", test.desc, test.err, err == nil)
			}
		}
	})

	t.Run("SignText", func(t *testing.T) {
		tests := []struct {
			name string
			key  *pgp.PrivateKey
			data []byte
			err  error
		}{
			{
				name: "valid key",
				key:  priv,
				data: data,
			},
			{
				name: "nil key",
				data: data,
				err:  crypto.ErrInvalidKey,
			},
		}

		for _, test := range tests {
			got, err := pgp.SignText(test.key, data)
			if !errors.Is(err, test.err) {
				t.Errorf("%s should return error %q got %q", test.name, test.err, err)
			} else if test.err == nil {
				gotData, gotSig, err := pgp.Split(got)
				if err != nil {
					t.Errorf("%s failed to split message: %s", test.name, err)
				} else if diff := cmp.Diff(string(test.data), string(gotData)); diff != "" {
					t.Error(test.name, diff)
				} else if test.err == nil && !pgp.Verify(pub, gotData, gotSig) {
					t.Error(test.name, "should pass verification")
				}
			}
		}
	})

	t.Run("Split", func(t *testing.T) {
		tests := []struct {
			name string
			in   []byte
			want []byte
			err  error
		}{
			{
				name: "valid message",
				in:   signed,
				want: data,
			},
			{
				name: "nil bytes",
				err:  pgp.ErrInvalidMessage,
			},
			{
				name: "unarmored data",
				in:   data,
				err:  pgp.ErrInvalidMessage,
			},
			{
				name: "missing signature",
				in:   []byte("-----BEGIN PGP SIGNED MESSAGE-----\r\nHash: SHA512\r\n\r\ndata\r\n"),
				err:  pgp.ErrInvalidMessage,
			},
			{
				name: "missing data",
				in:   append([]byte{'\n'}, sig...),
				err:  pgp.ErrInvalidMessage,
			},
		}

		for _, test := range tests {
			gotData, gotSig, err := pgp.Split(test.in)
			if !errors.Is(err, test.err) {
				t.Errorf("%s should return error %q got %q", test.name, test.err, err)
			} else if diff := cmp.Diff(string(test.want), string(gotData)); diff != "" {
				t.Error(test.name, diff)
			} else if test.err == nil && !pgp.Verify(pub, gotData, gotSig) {
				t.Error(test.name, "should pass verification")
			}
		}
	})

	t.Run("WrongKeyType", func(t *testing.T) {
		t.Run("MarshalPrivateKey", func(t *testing.T) {
			if _, err := pgp.MarshalPrivateKey(pub); err == nil {
				t.Errorf("should return error")
			}
		})

		t.Run("MarshalPublicKey", func(t *testing.T) {
			if _, err := pgp.MarshalPublicKey(priv); err == nil {
				t.Errorf("should return error")
			}
		})

		t.Run("Sign", func(t *testing.T) {
			if _, err := pgp.Sign(pub, data); err == nil {
				t.Errorf("should return error")
			}
		})
	})

	t.Run("LockedKey", func(t *testing.T) {
		privLocked, _ := priv.Lock([]byte("passphrase"))

		t.Run("Sign", func(t *testing.T) {
			if _, err := pgp.Sign(privLocked, data); err == nil {
				t.Errorf("should return error")
			}
		})

		t.Run("Verify", func(t *testing.T) {
			if pgp.Verify(privLocked, data, sig) {
				t.Errorf("should fail")
			}
		})
	})

	t.Run("GnuPG", func(t *testing.T) {
		if _, err := exec.LookPath("gpg"); err != nil {
			t.Skip("gpg not in path")
		}

		dir := t.TempDir()
		os.WriteFile(filepath.Join(dir, "key.asc"), pubBytes, 0o600)
		os.WriteFile(filepath.Join(dir, "data"), data, 0o600)
		os.WriteFile(filepath.Join(dir, "data.asc"), sig, 0o600)
		os.WriteFile(filepath.Join(dir, "signed"), signed, 0o600)

		baseArgs := []string{"--no-default-keyring", "--keyring", "keyring.gpg"}
		arguments := [][]string{
			{"--import", "key.asc"},  // Create keybox & import key.
			{"--verify", "data.asc"}, // Verify detached signature.
			{"--verify", "signed"},   // Verify signed message.
		}

		for _, a := range arguments {
			cmd := exec.Command("gpg", append(baseArgs, a...)...)
			cmd.Dir = dir
			out, err := cmd.CombinedOutput()
			if err != nil {
				t.Fatal(a, err, string(out))
			}
			t.Log(a, "\n"+string(out))
		}
	})
}
