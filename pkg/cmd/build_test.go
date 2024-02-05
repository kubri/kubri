package cmd_test

import (
	"bytes"
	"os"
	"strings"
	"testing"

	"github.com/abemedia/appcast/internal/test"
	"github.com/abemedia/appcast/pkg/cmd"
)

func TestBuild(t *testing.T) {
	tests := []struct {
		desc   string
		args   []string
		path   string
		config string
		want   string
		err    bool
	}{
		{
			desc: "no config file",
			args: []string{"build"},
			path: "test.yml",
			want: "Error: no config file found",
			err:  true,
		},
		{
			desc: "custom config file path",
			args: []string{"build", "-c", "test.yml"},
			path: "test.yml",
			want: "Error: no integrations configured",
			err:  true,
		},
		{
			desc:   "apk",
			args:   []string{"build"},
			path:   "appcast.yml",
			config: "apk: {}",
			want:   "Completed publishing APK packages.",
		},
		{
			desc:   "appinstaller",
			args:   []string{"build"},
			path:   "appcast.yml",
			config: "appinstaller: {}",
			want:   "Completed publishing App Installer packages.",
		},
		{
			desc:   "apt",
			args:   []string{"build"},
			path:   "appcast.yml",
			config: "apt: {}",
			want:   "Completed publishing APT packages.",
		},
		{
			desc:   "yum",
			args:   []string{"build"},
			path:   "appcast.yml",
			config: "yum: {}",
			want:   "Completed publishing YUM packages.",
		},
		{
			desc:   "sparkle",
			args:   []string{"build"},
			path:   "appcast.yml",
			config: "sparkle: {}",
			want:   "Completed publishing Sparkle packages.",
		},
	}

	baseConfig := `
		source:
			type: file
			path: ` + t.TempDir() + `
		target:
			type: file
			path: ` + t.TempDir()

	for _, tc := range tests {
		os.Chdir(t.TempDir())
		os.WriteFile(tc.path, test.JoinYAML(tc.config, baseConfig), os.ModePerm)

		var out bytes.Buffer
		err := cmd.Execute("", cmd.WithArgs(tc.args...), cmd.WithStderr(&out), cmd.WithStdout(&out))
		if tc.err != (err != nil) || !strings.Contains(out.String(), tc.want) {
			t.Errorf("%s should return %q:\n%s", tc.desc, tc.want, &out)
		}
	}
}
