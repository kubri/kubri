package cmd_test

import (
	"bytes"
	"io"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/kubri/kubri/pkg/cmd"
	"github.com/kubri/kubri/pkg/crypto/dsa"
	"github.com/kubri/kubri/pkg/crypto/ed25519"
	"github.com/kubri/kubri/pkg/crypto/pgp"
	"github.com/kubri/kubri/pkg/crypto/rsa"
	"github.com/kubri/kubri/pkg/secret"
)

func TestKeysImportCmd(t *testing.T) {
	dsaKey, _ := dsa.NewPrivateKey()
	dsaBytes, _ := dsa.MarshalPrivateKey(dsaKey)
	dsaPath := filepath.Join(t.TempDir(), "test")
	os.WriteFile(dsaPath, dsaBytes, os.ModePerm)

	ed25519Key, _ := ed25519.NewPrivateKey()
	ed25519Bytes, _ := ed25519.MarshalPrivateKey(ed25519Key)
	ed25519Path := filepath.Join(t.TempDir(), "test")
	os.WriteFile(ed25519Path, ed25519Bytes, os.ModePerm)

	ed25519PEMBytes, _ := ed25519.MarshalPrivateKeyPEM(ed25519Key)
	ed25519PEMPath := filepath.Join(t.TempDir(), "test")
	os.WriteFile(ed25519PEMPath, ed25519PEMBytes, os.ModePerm)

	pgpKey, _ := pgp.NewPrivateKey("test", "test@example.com")
	pgpBytes, _ := pgp.MarshalPrivateKey(pgpKey)
	pgpPath := filepath.Join(t.TempDir(), "test")
	os.WriteFile(pgpPath, pgpBytes, os.ModePerm)

	rsaKey, _ := rsa.NewPrivateKey()
	rsaBytes, _ := rsa.MarshalPrivateKey(rsaKey)
	rsaPath := filepath.Join(t.TempDir(), "test")
	os.WriteFile(rsaPath, rsaBytes, os.ModePerm)

	tests := []struct {
		args []string
		want []byte
	}{
		{
			args: []string{"keys", "import", "dsa", dsaPath},
			want: dsaBytes,
		},
		{
			args: []string{"keys", "import", "ed25519", ed25519Path},
			want: ed25519Bytes,
		},
		{
			args: []string{"keys", "import", "pgp", pgpPath},
			want: pgpBytes,
		},
		{
			args: []string{"keys", "import", "rsa", rsaPath},
			want: rsaBytes,
		},
		{
			args: []string{"keys", "import", "dsa", dsaPath, "--force"},
			want: dsaBytes,
		},
		{
			args: []string{"keys", "import", "ed25519", ed25519Path, "--force"},
			want: ed25519Bytes,
		},
		{
			args: []string{"keys", "import", "ed25519", ed25519PEMPath, "--force"},
			want: ed25519Bytes,
		},
		{
			args: []string{"keys", "import", "pgp", pgpPath, "--force"},
			want: pgpBytes,
		},
		{
			args: []string{"keys", "import", "rsa", rsaPath, "--force"},
			want: rsaBytes,
		},
	}

	t.Setenv("KUBRI_PATH", t.TempDir())

	for _, test := range tests {
		err := cmd.Execute("", cmd.WithArgs(test.args...), cmd.WithStdout(io.Discard))
		if err != nil {
			t.Errorf("%s: %s", test, err)
		} else if b, _ := secret.Get(test.args[2] + "_key"); !bytes.Equal(test.want, b) {
			t.Errorf("%s should be equal", test)
		}
	}
}

func TestKeysImportCmdErrors(t *testing.T) {
	path := filepath.Join(t.TempDir(), "test")
	os.WriteFile(path, []byte("test"), os.ModePerm)

	dsaPath := filepath.Join(t.TempDir(), "test")
	dsaKey, _ := dsa.NewPrivateKey()
	dsaBytes, _ := dsa.MarshalPrivateKey(dsaKey)
	os.WriteFile(dsaPath, dsaBytes, os.ModePerm)

	ed25519Path := filepath.Join(t.TempDir(), "test")
	ed25519Key, _ := ed25519.NewPrivateKey()
	ed25519Bytes, _ := ed25519.MarshalPrivateKey(ed25519Key)
	os.WriteFile(ed25519Path, ed25519Bytes, os.ModePerm)

	pgpPath := filepath.Join(t.TempDir(), "test")
	pgpKey, _ := pgp.NewPrivateKey("test", "test@example.com")
	pgpBytes, _ := pgp.MarshalPrivateKey(pgpKey)
	os.WriteFile(pgpPath, pgpBytes, os.ModePerm)

	rsaPath := filepath.Join(t.TempDir(), "test")
	rsaKey, _ := rsa.NewPrivateKey()
	rsaBytes, _ := rsa.MarshalPrivateKey(rsaKey)
	os.WriteFile(rsaPath, rsaBytes, os.ModePerm)

	tests := []struct {
		args []string
		want string
	}{
		{
			args: []string{"keys", "import"},
			want: "accepts 2 arg(s), received 0",
		},
		{
			args: []string{"keys", "import", "dsa"},
			want: "accepts 2 arg(s), received 1",
		},
		{
			args: []string{"keys", "import", "foo", "bar"},
			want: "invalid argument",
		},
		{
			args: []string{"keys", "import", "dsa", dsaPath},
			want: "key already exists",
		},
		{
			args: []string{"keys", "import", "ed25519", ed25519Path},
			want: "key already exists",
		},
		{
			args: []string{"keys", "import", "pgp", pgpPath},
			want: "key already exists",
		},
		{
			args: []string{"keys", "import", "rsa", rsaPath},
			want: "key already exists",
		},
		{
			args: []string{"keys", "import", "dsa", path},
			want: "invalid key",
		},
		{
			args: []string{"keys", "import", "ed25519", path},
			want: "invalid key",
		},
		{
			args: []string{"keys", "import", "pgp", path},
			want: "invalid key",
		},
		{
			args: []string{"keys", "import", "rsa", path},
			want: "invalid key",
		},
		{
			args: []string{"keys", "import", "dsa", "foo", "--force"},
			want: "no such file or directory",
		},
	}

	t.Setenv("KUBRI_PATH", t.TempDir())
	secret.Put("dsa_key", nil)
	secret.Put("ed25519_key", nil)
	secret.Put("pgp_key", nil)
	secret.Put("rsa_key", nil)

	for _, test := range tests {
		var stderr bytes.Buffer
		err := cmd.Execute("", cmd.WithArgs(test.args...), cmd.WithStderr(&stderr), cmd.WithStdout(io.Discard))
		if err == nil || !strings.Contains(stderr.String(), test.want) {
			t.Errorf("%s should fail with %q:\n%s", test.args, test.want, err)
		}
	}
}
