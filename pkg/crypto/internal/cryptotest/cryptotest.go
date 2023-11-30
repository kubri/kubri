package cryptotest

import (
	"bytes"
	"errors"
	"os"
	"os/exec"
	"path/filepath"
	"testing"

	"github.com/abemedia/appcast/pkg/crypto"
	"github.com/google/go-cmp/cmp"
)

type Implementation[PrivateKey any, PublicKey any] struct {
	NewPrivateKey       func() (PrivateKey, error)
	MarshalPrivateKey   func(key PrivateKey) ([]byte, error)
	UnmarshalPrivateKey func(b []byte) (PrivateKey, error)
	Public              func(key PrivateKey) PublicKey
	MarshalPublicKey    func(key PublicKey) ([]byte, error)
	UnmarshalPublicKey  func(b []byte) (PublicKey, error)
	Sign                func(key PrivateKey, data []byte) ([]byte, error)
	Verify              func(key PublicKey, data, sig []byte) bool
}

type options struct {
	cmp         cmp.Options
	opensslArgs []string
}

type Option func(*options)

func WithCmpOptions(opt ...cmp.Option) Option {
	return func(o *options) { o.cmp = append(o.cmp, opt...) }
}

func WithOpenSSLTest(arg ...string) Option {
	return func(o *options) { o.opensslArgs = arg }
}

//nolint:funlen,gocognit,maintidx
func Test[PrivateKey, PublicKey any](t *testing.T, i Implementation[PrivateKey, PublicKey], opts ...Option) {
	var opt options
	for _, o := range opts {
		o(&opt)
	}

	priv, err := i.NewPrivateKey()
	if err != nil {
		t.Fatal(err)
	}

	privBytes, err := i.MarshalPrivateKey(priv)
	if err != nil {
		t.Fatal(err)
	}

	pub := i.Public(priv)

	pubBytes, err := i.MarshalPublicKey(pub)
	if err != nil {
		t.Fatal(err)
	}

	data := []byte("foo\nbar\nbaz\n")

	sig, err := i.Sign(priv, data)
	if err != nil {
		t.Fatal(err)
	}

	t.Run("MarshalPrivateKey", func(t *testing.T) {
		tests := []struct {
			name string
			in   PrivateKey
			want []byte
			err  error
		}{
			{
				name: "valid key",
				in:   priv,
				want: privBytes,
			},
			{
				name: "nil key",
				err:  crypto.ErrInvalidKey,
			},
		}

		for _, test := range tests {
			got, err := i.MarshalPrivateKey(test.in)
			if !errors.Is(err, test.err) {
				t.Error(test.name, "should return error", test.err, "got", err)
			} else if diff := cmp.Diff(string(test.want), string(got), opt.cmp); diff != "" {
				t.Error(test.name, diff)
			}
		}
	})

	t.Run("UnmarshalPrivateKey", func(t *testing.T) {
		tests := []struct {
			name string
			in   []byte
			want PrivateKey
			err  error
		}{
			{
				name: "valid key",
				in:   privBytes,
				want: priv,
			},
			{
				name: "nil bytes",
				err:  crypto.ErrInvalidKey,
			},
			{
				name: "non-key data",
				in:   data,
				err:  crypto.ErrInvalidKey,
			},
			{
				name: "public key",
				in:   pubBytes,
				err:  crypto.ErrInvalidKey,
			},
		}

		for _, test := range tests {
			got, err := i.UnmarshalPrivateKey(test.in)
			if !errors.Is(err, test.err) {
				t.Errorf("%s should return error %q got %q", test.name, test.err, err)
			} else if diff := cmp.Diff(test.want, got, opt.cmp); diff != "" {
				t.Error(test.name, diff)
			}
		}
	})

	t.Run("MarshalPublicKey", func(t *testing.T) {
		tests := []struct {
			name string
			in   PublicKey
			want []byte
			err  error
		}{
			{
				name: "valid key",
				in:   pub,
				want: pubBytes,
			},
			{
				name: "nil key",
				err:  crypto.ErrInvalidKey,
			},
		}

		for _, test := range tests {
			got, err := i.MarshalPublicKey(test.in)
			if !errors.Is(err, test.err) {
				t.Errorf("%s should return error %q got %q", test.name, test.err, err)
			} else if diff := cmp.Diff(string(test.want), string(got), opt.cmp); diff != "" {
				t.Error(test.name, diff)
			}
		}
	})

	t.Run("UnmarshalPublicKey", func(t *testing.T) {
		tests := []struct {
			name string
			in   []byte
			want PublicKey
			err  error
		}{
			{
				name: "valid key",
				in:   pubBytes,
				want: pub,
			},
			{
				name: "nil bytes",
				err:  crypto.ErrInvalidKey,
			},
			{
				name: "non-key data",
				in:   data,
				err:  crypto.ErrInvalidKey,
			},
			{
				name: "private key",
				in:   privBytes,
				err:  crypto.ErrInvalidKey,
			},
		}

		for _, test := range tests {
			got, err := i.UnmarshalPublicKey(test.in)
			if !errors.Is(err, test.err) {
				t.Errorf("%s should return error %q got %q", test.name, test.err, err)
			} else if diff := cmp.Diff(test.want, got, opt.cmp); diff != "" {
				t.Errorf("%s\n%s", test.name, diff)
			}
		}
	})

	t.Run("Sign", func(t *testing.T) {
		tests := []struct {
			name string
			key  PrivateKey
			data []byte
			err  error
		}{
			{
				name: "nil key",
				data: data,
			},
		}

		for _, test := range tests {
			_, err := i.Sign(test.key, data)
			if err == nil {
				t.Error(test.name, "should error")
			}
		}
	})

	t.Run("Verify", func(t *testing.T) {
		wrongPriv, _ := i.NewPrivateKey()
		wrongPub := i.Public(wrongPriv)
		wrongSig, _ := i.Sign(priv, []byte("wrong data"))

		tests := []struct {
			name string
			key  PublicKey
			data []byte
			sig  []byte
			want bool
		}{
			{
				name: "valid key",
				key:  pub,
				data: data,
				sig:  sig,
				want: true,
			},
			{
				name: "nil key",
				data: data,
				sig:  sig,
			},
			{
				name: "wrong key",
				key:  wrongPub,
				data: data,
				sig:  sig,
			},
			{
				name: "nil data",
				key:  pub,
				sig:  sig,
			},
			{
				name: "nil signature",
				key:  pub,
				data: data,
			},
			{
				name: "wrong signature",
				key:  pub,
				data: data,
				sig:  wrongSig,
			},
		}

		for _, test := range tests {
			ok := i.Verify(test.key, test.data, test.sig)
			if ok != test.want {
				t.Errorf("%s should return %t got %t", test.name, test.want, ok)
			}
		}
	})

	if opt.opensslArgs != nil {
		t.Run("OpenSSL", func(t *testing.T) {
			dir := t.TempDir()
			_ = os.WriteFile(filepath.Join(dir, "public.pem"), pubBytes, 0o600)
			_ = os.WriteFile(filepath.Join(dir, "data.txt"), data, 0o600)
			_ = os.WriteFile(filepath.Join(dir, "data.txt.sig"), sig, 0o600)

			cmd := exec.Command("openssl", opt.opensslArgs...)
			cmd.Dir = dir

			out, err := cmd.CombinedOutput()
			t.Log(string(bytes.TrimSpace(out)))
			if err != nil {
				t.Fatal(err)
			}
		})
	}
}
