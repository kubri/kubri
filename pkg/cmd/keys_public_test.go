package cmd_test

import (
	"bytes"
	"io"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/abemedia/appcast/pkg/cmd"
	"github.com/abemedia/appcast/pkg/crypto/dsa"
	"github.com/abemedia/appcast/pkg/crypto/ed25519"
	"github.com/abemedia/appcast/pkg/crypto/pgp"
	"github.com/abemedia/appcast/pkg/crypto/rsa"
)

func TestKeysPublicCmd(t *testing.T) {
	dir := t.TempDir()
	t.Setenv("APPCAST_PATH", dir)

	{
		var stdout bytes.Buffer
		key, _ := dsa.NewPrivateKey()
		priv, _ := dsa.MarshalPrivateKey(key)
		os.WriteFile(filepath.Join(dir, "dsa_key"), priv, 0o600)
		pub, _ := dsa.MarshalPublicKey(dsa.Public(key))

		err := cmd.Execute("", cmd.WithArgs("keys", "public", "dsa"), cmd.WithStdout(&stdout))
		if err != nil {
			t.Error(err)
		}
		if !bytes.Equal(stdout.Bytes(), pub) {
			t.Error("dsa public keys should be equal")
		}
	}

	{
		var stdout bytes.Buffer
		key, _ := ed25519.NewPrivateKey()
		priv, _ := ed25519.MarshalPrivateKey(key)
		os.WriteFile(filepath.Join(dir, "ed25519_key"), priv, 0o600)
		pub, _ := ed25519.MarshalPublicKey(ed25519.Public(key))

		err := cmd.Execute("", cmd.WithArgs("keys", "public", "ed25519"), cmd.WithStdout(&stdout))
		if err != nil {
			t.Error(err)
		}
		if !bytes.Equal(stdout.Bytes(), pub) {
			t.Error("ed25519 public keys should be equal")
		}
	}

	{
		var stdout bytes.Buffer
		key, _ := pgp.NewPrivateKey("test", "test@example.com")
		priv, _ := pgp.MarshalPrivateKey(key)
		os.WriteFile(filepath.Join(dir, "pgp_key"), priv, 0o600)
		pub, _ := pgp.MarshalPublicKey(pgp.Public(key))

		err := cmd.Execute("", cmd.WithArgs("keys", "public", "pgp"), cmd.WithStdout(&stdout))
		if err != nil {
			t.Error(err)
		}
		if !bytes.Equal(stdout.Bytes(), pub) {
			t.Error("pgp public keys should be equal")
		}
	}

	{
		var stdout bytes.Buffer
		key, _ := rsa.NewPrivateKey()
		priv, _ := rsa.MarshalPrivateKey(key)
		os.WriteFile(filepath.Join(dir, "rsa_key"), priv, 0o600)
		pub, _ := rsa.MarshalPublicKey(rsa.Public(key))

		err := cmd.Execute("", cmd.WithArgs("keys", "public", "rsa"), cmd.WithStdout(&stdout))
		if err != nil {
			t.Error(err)
		}
		if !bytes.Equal(stdout.Bytes(), pub) {
			t.Error("rsa public keys should be equal")
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
		{
			args: []string{"keys", "public", "rsa"},
			want: "key not found",
		},
	}

	t.Setenv("APPCAST_PATH", t.TempDir())

	for _, test := range tests {
		var stderr bytes.Buffer
		err := cmd.Execute("", cmd.WithArgs(test.args...), cmd.WithStderr(&stderr), cmd.WithStdout(io.Discard))
		if err == nil || !strings.Contains(stderr.String(), test.want) {
			t.Errorf("%s should fail with %q:\n%s", test.args, test.want, &stderr)
		}
	}
}
