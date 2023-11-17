package cmd_test

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/abemedia/appcast/pkg/cmd"
	"github.com/abemedia/appcast/pkg/crypto/dsa"
	"github.com/abemedia/appcast/pkg/crypto/ed25519"
	"github.com/abemedia/appcast/pkg/crypto/pgp"
)

func TestKeysPublicCmd(t *testing.T) {
	stdout := capture(t, os.Stdout)
	t.Cleanup(func() { t.Log(stdout.String()) })
	dir := t.TempDir()
	t.Setenv("APPCAST_PATH", dir)

	{
		key, _ := dsa.NewPrivateKey()
		priv, _ := dsa.MarshalPrivateKey(key)
		os.WriteFile(filepath.Join(dir, "dsa_key"), priv, 0o600)
		pub, _ := dsa.MarshalPublicKey(dsa.Public(key))

		err := cmd.Execute("", []string{"keys", "public", "dsa"})
		if err != nil {
			t.Error(err)
		}
		if !bytes.Equal(stdout.Bytes(), pub) {
			t.Error("dsa public keys should be equal")
		}
	}

	stdout.Reset()

	{
		key, _ := ed25519.NewPrivateKey()
		priv, _ := ed25519.MarshalPrivateKey(key)
		os.WriteFile(filepath.Join(dir, "ed25519_key"), priv, 0o600)
		pub, _ := ed25519.MarshalPublicKey(ed25519.Public(key))

		err := cmd.Execute("", []string{"keys", "public", "ed25519"})
		if err != nil {
			t.Error(err)
		}
		if !bytes.Equal(stdout.Bytes(), pub) {
			t.Error("ed25519 public keys should be equal")
		}
	}

	stdout.Reset()

	{
		key, _ := pgp.NewPrivateKey("test", "test@example.com")
		priv, _ := pgp.MarshalPrivateKey(key)
		os.WriteFile(filepath.Join(dir, "pgp_key"), priv, 0o600)
		pub, _ := pgp.MarshalPublicKey(pgp.Public(key))

		err := cmd.Execute("", []string{"keys", "public", "pgp"})
		if err != nil {
			t.Error(err)
		}
		if !bytes.Equal(stdout.Bytes(), pub) {
			t.Error("pgp public keys should be equal")
		}
	}
}

func TestKeysPublicCmdErrors(t *testing.T) {
	tests := []struct {
		args []string
		want string
	}{
		{
			args: []string{"keys", "public"},
			want: "accepts 1 arg(s), received 0",
		},
		{
			args: []string{"keys", "public", "foo"},
			want: "invalid argument",
		},
		{
			args: []string{"keys", "public", "dsa"},
			want: "key not found",
		},
		{
			args: []string{"keys", "public", "ed25519"},
			want: "key not found",
		},
	}

	stderr := capture(t, os.Stderr)
	t.Setenv("APPCAST_PATH", t.TempDir())

	for _, test := range tests {
		err := cmd.Execute("", test.args)
		if err == nil || !strings.Contains(stderr.String(), test.want) {
			t.Errorf("%s should fail with %q:\n%s", test.args, test.want, stderr)
		}
		stderr.Reset()
	}
}
