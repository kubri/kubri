package cmd_test

import (
	"bytes"
	"io"
	"os"
	"strings"
	"testing"

	"github.com/abemedia/appcast/pkg/cmd"
)

func TestBuild(t *testing.T) {
	tests := []struct {
		args []string
		path string
		want string
		err  error
	}{
		{
			args: []string{"build"},
			path: "test.yml",
			want: "no config file found",
		},
		{
			args: []string{"build", "-c", "test.yml"},
			path: "test.yml",
			want: "no integrations configured",
		},
		{
			args: []string{"build"},
			path: "appcast.yml",
			want: "no integrations configured",
		},
	}

	config := `
source:
  type: file
  path: ` + t.TempDir() + `
target:
  type: file
  path: ` + t.TempDir()

	for _, test := range tests {
		os.Chdir(t.TempDir())
		os.WriteFile(test.path, []byte(config), os.ModePerm)

		var stderr bytes.Buffer
		err := cmd.Execute("", cmd.WithArgs(test.args...), cmd.WithStderr(&stderr), cmd.WithStdout(io.Discard))
		if err == nil || !strings.Contains(stderr.String(), test.want) {
			t.Errorf("%s should fail with %q:\n%s", test.args, test.want, &stderr)
		}
	}
}
