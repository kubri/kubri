package cmd_test

import (
	"bytes"
	"io"
	"strings"
	"testing"

	"github.com/kubri/kubri/pkg/cmd"
	"github.com/kubri/kubri/pkg/crypto/dsa"
	"github.com/kubri/kubri/pkg/crypto/ed25519"
	"github.com/kubri/kubri/pkg/crypto/pgp"
	"github.com/kubri/kubri/pkg/crypto/rsa"
	"github.com/kubri/kubri/pkg/secret"
)

func TestKeysExportCmd(t *testing.T) {
	t.Setenv("KUBRI_PATH", t.TempDir())

	dsaKey, _ := dsa.NewPrivateKey()
	dsaBytes, _ := dsa.MarshalPrivateKey(dsaKey)
	secret.Put("dsa_key", dsaBytes)

	ed25519Key, _ := ed25519.NewPrivateKey()
	ed25519Bytes, _ := ed25519.MarshalPrivateKey(ed25519Key)
	secret.Put("ed25519_key", ed25519Bytes)

	pgpKey, _ := pgp.NewPrivateKey("test", "test@example.com")
	pgpBytes, _ := pgp.MarshalPrivateKey(pgpKey)
	secret.Put("pgp_key", pgpBytes)

	rsaKey, _ := rsa.NewPrivateKey()
	rsaBytes, _ := rsa.MarshalPrivateKey(rsaKey)
	secret.Put("rsa_key", rsaBytes)

	tests := [][]string{
		{"keys", "export", "dsa"},
		{"keys", "export", "ed25519"},
		{"keys", "export", "pgp"},
		{"keys", "export", "rsa"},
	}

	for _, test := range tests {
		want, _ := secret.Get(test[2] + "_key")
		err := cmd.Execute("", cmd.WithArgs(test...), cmd.WithStdout(io.Discard))
		if err != nil {
			t.Errorf("%s: %s", test, err)
		} else if b, _ := secret.Get(test[2] + "_key"); !bytes.Equal(want, b) {
			t.Errorf("%s should be equal", test)
		}
	}
}

func TestKeysExportCmdErrors(t *testing.T) {
	tests := []struct {
		args []string
		want string
	}{
		{
			args: []string{"keys", "export"},
			want: "accepts 1 arg(s), received 0",
		},
		{
			args: []string{"keys", "export", "foo"},
			want: "invalid argument",
		},
		{
			args: []string{"keys", "export", "foo", "bar"},
			want: "accepts 1 arg(s), received 2",
		},
		{
			args: []string{"keys", "export", "dsa", "rsa"},
			want: "accepts 1 arg(s), received 2",
		},
		{
			args: []string{"keys", "export", "dsa"},
			want: "key not found",
		},
		{
			args: []string{"keys", "export", "ed25519"},
			want: "key not found",
		},
		{
			args: []string{"keys", "export", "pgp"},
			want: "key not found",
		},
		{
			args: []string{"keys", "export", "rsa"},
			want: "key not found",
		},
	}

	t.Setenv("KUBRI_PATH", t.TempDir())

	for _, test := range tests {
		var stderr bytes.Buffer
		err := cmd.Execute("", cmd.WithArgs(test.args...), cmd.WithStderr(&stderr), cmd.WithStdout(io.Discard))
		if err == nil || !strings.Contains(stderr.String(), test.want) {
			t.Errorf("%s should fail with %q:\n%s", test.args, test.want, err)
		}
	}
}
