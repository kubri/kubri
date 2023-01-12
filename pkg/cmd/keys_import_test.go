package cmd_test

import (
	"bytes"
	"crypto/rand"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/abemedia/appcast/pkg/cmd"
	"github.com/abemedia/appcast/pkg/secret"
)

func TestKeysImportCmd(t *testing.T) {
	tests := [][]string{
		{"keys", "import", "dsa", filepath.Join(t.TempDir(), "test")},
		{"keys", "import", "ed25519", filepath.Join(t.TempDir(), "test")},
		{"keys", "import", "dsa", filepath.Join(t.TempDir(), "test"), "--force"},
		{"keys", "import", "ed25519", filepath.Join(t.TempDir(), "test"), "--force"},
	}

	capture(t, os.Stderr)
	t.Setenv("APPCAST_PATH", t.TempDir())

	for _, test := range tests {
		want := make([]byte, 10)
		rand.Read(want)
		os.WriteFile(test[3], want, os.ModePerm)

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
			args: []string{"keys", "import", "dsa", path},
			want: "key already exists",
		},
		{
			args: []string{"keys", "import", "ed25519", path},
			want: "key already exists",
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

	for _, test := range tests {
		err := cmd.Execute("", test.args)
		if err == nil || !strings.Contains(stderr.String(), test.want) {
			t.Errorf("%s should fail with %q:\n%s", test.args, test.want, err)
		}
		stderr.Reset()
	}
}
