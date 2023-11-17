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
	"github.com/abemedia/appcast/pkg/secret"
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

	pgpKey, _ := pgp.NewPrivateKey("test", "test@example.com")
	pgpBytes, _ := pgp.MarshalPrivateKey(pgpKey)
	pgpPath := filepath.Join(t.TempDir(), "test")
	os.WriteFile(pgpPath, pgpBytes, os.ModePerm)

	tests := [][]string{
		{"keys", "import", "dsa", dsaPath},
		{"keys", "import", "ed25519", ed25519Path},
		{"keys", "import", "pgp", pgpPath},
		{"keys", "import", "dsa", dsaPath, "--force"},
		{"keys", "import", "ed25519", ed25519Path, "--force"},
		{"keys", "import", "pgp", pgpPath, "--force"},
	}

	capture(t, os.Stderr)
	t.Setenv("APPCAST_PATH", t.TempDir())

	for _, test := range tests {
		want, _ := os.ReadFile(test[3])
		err := cmd.Execute("", test)
		if err != nil {
			t.Errorf("%s: %s", test, err)
		} else if b, _ := secret.Get(test[2] + "_key"); !bytes.Equal(want, b) {
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
			args: []string{"keys", "import", "dsa", "foo", "--force"},
			want: "no such file or directory",
		},
	}

	stderr := capture(t, os.Stderr)
	t.Setenv("APPCAST_PATH", t.TempDir())
	secret.Put("dsa_key", nil)
	secret.Put("ed25519_key", nil)
	secret.Put("pgp_key", nil)

	for _, test := range tests {
		err := cmd.Execute("", test.args)
		if err == nil || !strings.Contains(stderr.String(), test.want) {
			t.Errorf("%s should fail with %q:\n%s", test.args, test.want, err)
		}
		stderr.Reset()
	}
}
