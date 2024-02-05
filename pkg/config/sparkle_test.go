package config_test

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/MakeNowJust/heredoc/v2"
	"github.com/abemedia/appcast/integrations/sparkle"
	"github.com/abemedia/appcast/pkg/config"
	"github.com/abemedia/appcast/pkg/crypto"
	"github.com/abemedia/appcast/pkg/crypto/dsa"
	"github.com/abemedia/appcast/pkg/crypto/ed25519"
	"github.com/abemedia/appcast/pkg/secret"
	source "github.com/abemedia/appcast/source/file"
	target "github.com/abemedia/appcast/target/file"
	"github.com/google/go-cmp/cmp"
)

func TestSparkle(t *testing.T) {
	dir := t.TempDir()
	src, _ := source.New(source.Config{Path: dir})
	tgt, _ := target.New(target.Config{Path: dir})
	dsaKey, _ := dsa.NewPrivateKey()
	dsaBytes, _ := dsa.MarshalPrivateKey(dsaKey)
	edKey, _ := ed25519.NewPrivateKey()
	edBytes, _ := ed25519.MarshalPrivateKey(edKey)

	runTest(t, []testCase{
		{
			desc: "disabled",
			in: `
				source:
					type: file
					path: ` + dir + `
				target:
					type: file
					path: ` + dir + `
				sparkle:
					disabled: true
			`,
			want: &config.Config{},
		},
		{
			desc: "defaults",
			in: `
				title: title
				description: description
				source:
					type: file
					path: ` + dir + `
				target:
					type: file
					path: ` + dir + `
				sparkle: {}
			`,
			want: &config.Config{
				Sparkle: &sparkle.Config{
					Title:       "title",
					Description: "description",
					FileName:    "appcast.xml",
					Settings:    []sparkle.Rule{},
					Source:      src,
					Target:      tgt.Sub("sparkle"),
				},
			},
		},
		{
			desc: "full",
			in: `
				title: ignore
				description: ignore
				version: latest
				prerelease: true
				upload-packages: true
				source:
					type: file
					path: ` + dir + `
				target:
					type: file
					path: ` + dir + `
				sparkle:
					title: title
					description: description
					folder: test
					filename: test.xml
					params:
						- os: windows
							installer-arguments: /passive
						- os: macos
							minimum-system-version: '10.13.0'
						- version: '1.0.0'
							critical-update: true
						- version: '> 1.0.0'
							critical-update-below-version: '1.0.0'
							minimum-autoupdate-version: '1.0.0'
						- version: '1.1.0'
							ignore-skipped-upgrades-below-version: '1.1.0'
			`,
			hook: func() {
				secret.Put("dsa_key", dsaBytes)
				secret.Put("ed25519_key", edBytes)
			},
			want: &config.Config{
				Sparkle: &sparkle.Config{
					Title:       "title",
					Description: "description",
					FileName:    "test.xml",
					DSAKey:      dsaKey,
					Ed25519Key:  edKey,
					Settings: []sparkle.Rule{
						{
							OS: sparkle.Windows,
							Settings: &sparkle.Settings{
								InstallerArguments: "/passive",
							},
						},
						{
							OS: sparkle.MacOS,
							Settings: &sparkle.Settings{
								MinimumSystemVersion: "10.13.0",
							},
						},
						{
							Version: "1.0.0",
							Settings: &sparkle.Settings{
								CriticalUpdate: true,
							},
						},
						{
							Version: "> 1.0.0",
							Settings: &sparkle.Settings{
								CriticalUpdateBelowVersion: "1.0.0",
								MinimumAutoupdateVersion:   "1.0.0",
							},
						},
						{
							Version: "1.1.0",
							Settings: &sparkle.Settings{
								IgnoreSkippedUpgradesBelowVersion: "1.1.0",
							},
						},
					},
					Source:         src,
					Target:         tgt.Sub("test"),
					Version:        "latest",
					Prerelease:     true,
					UploadPackages: true,
				},
			},
		},
		{
			desc: "invalid dsa key",
			in: `
				source:
					type: file
					path: ` + dir + `
				target:
					type: file
					path: ` + dir + `
				sparkle: {}
			`,
			hook: func() { secret.Put("dsa_key", []byte("nope")) },
			err:  crypto.ErrInvalidKey,
		},
		{
			desc: "invalid ed25519 key",
			in: `
				source:
					type: file
					path: ` + dir + `
				target:
					type: file
					path: ` + dir + `
				sparkle: {}
			`,
			hook: func() { secret.Put("ed25519_key", []byte("nope")) },
			err:  crypto.ErrInvalidKey,
		},
	})
}

func TestSparkleDetectOS(t *testing.T) {
	c := `
		source:
			type: file
			path: ` + t.TempDir() + `
		target:
			type: file
			path: ` + t.TempDir() + `
		sparkle:
			detect-os:
				macos: '*.dmg'
				windows: '*.exe'
	`

	tests := []struct {
		in   string
		want sparkle.OS
	}{
		{
			in:   "test.dmg",
			want: sparkle.MacOS,
		},
		{
			in:   "test.exe",
			want: sparkle.Windows,
		},
		{
			in:   "unknown",
			want: sparkle.Unknown,
		},
	}

	t.Setenv("APPCAST_PATH", t.TempDir())
	path := filepath.Join(t.TempDir(), "appcast.yml")
	os.WriteFile(path, []byte(strings.ReplaceAll(heredoc.Doc(c), "\t", "  ")), os.ModePerm)

	p, err := config.Load(path)
	if err != nil {
		t.Fatal(err)
	}

	for i, test := range tests {
		got := p.Sparkle.DetectOS(test.in)
		if diff := cmp.Diff(test.want, got); diff != "" {
			t.Error(i, diff)
		}
	}
}
