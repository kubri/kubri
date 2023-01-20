package pipe_test

import (
	"encoding/pem"
	"os"
	"path/filepath"
	"reflect"
	"testing"

	"github.com/abemedia/appcast/integrations/appinstaller"
	"github.com/abemedia/appcast/integrations/apt"
	"github.com/abemedia/appcast/integrations/sparkle"
	"github.com/abemedia/appcast/pkg/crypto/dsa"
	"github.com/abemedia/appcast/pkg/crypto/ed25519"
	"github.com/abemedia/appcast/pkg/pipe"
	fileSource "github.com/abemedia/appcast/source/blob/file"
	fileTarget "github.com/abemedia/appcast/target/file"
	"github.com/google/go-cmp/cmp"
)

func TestPipe(t *testing.T) {
	dir := t.TempDir()

	src, err := fileSource.New(fileSource.Config{Path: dir})
	if err != nil {
		t.Fatal(err)
	}

	tgt, err := fileTarget.New(fileTarget.Config{Path: dir})
	if err != nil {
		t.Fatal(err)
	}

	dsaKey, _ := dsa.NewPrivateKey()
	edKey, _ := ed25519.NewPrivateKey()

	tests := []struct {
		in   string
		want *pipe.Pipe
	}{
		{
			in: `
title: test
description: test
source:
  type: file
  path: ` + dir + `
target:
  type: file
  path: ` + dir + `
appinstaller:
  disabled: true
apt:
  disabled: true
sparkle:
  disabled: true
`,
			want: &pipe.Pipe{},
		},
		{
			in: `
title: test
description: test
source:
  type: file
  path: ` + dir + `
target:
  type: file
  path: ` + dir + `
`,
			want: &pipe.Pipe{
				Appinstaller: &appinstaller.Config{
					Source: src,
					Target: tgt.Sub("appinstaller"),
				},
				Apt: &apt.Config{
					Source: src,
					Target: tgt.Sub("apt"),
				},
				Sparkle: &sparkle.Config{
					Title:       "test",
					Description: "test",
					Source:      src,
					Target:      tgt.Sub("sparkle"),
					FileName:    "sparkle.xml",
					Settings:    []sparkle.Rule{},
				},
			},
		},
		{
			in: `
title: test
description: test
version: latest
source:
  type: file
  path: ` + dir + `
target:
  type: file
  path: ` + dir + `
appinstaller:
  folder: .
  hours-between-update-checks: 12
  update-blocks-activation: true
  show-prompt: true
  automatic-background-task: true
  force-update-from-any-version: true	
apt:
  folder: .
sparkle:
  folder: .
  filename: updates.xml
  title: foo
  description: bar
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
			want: &pipe.Pipe{
				Appinstaller: &appinstaller.Config{
					Source:                    src,
					Target:                    tgt,
					Version:                   "latest",
					HoursBetweenUpdateChecks:  12,
					UpdateBlocksActivation:    true,
					ShowPrompt:                true,
					AutomaticBackgroundTask:   true,
					ForceUpdateFromAnyVersion: true,
				},
				Apt: &apt.Config{
					Source:  src,
					Target:  tgt,
					Version: "latest",
				},
				Sparkle: &sparkle.Config{
					Title:       "foo",
					Description: "bar",
					Source:      src,
					Target:      tgt,
					FileName:    "updates.xml",
					Version:     "latest",
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
				},
			},
		},
	}

	for i, test := range tests {
		dir := t.TempDir()
		t.Setenv("APPCAST_PATH", dir)

		if test.want.Sparkle != nil {
			if test.want.Sparkle.DSAKey != nil {
				b, _ := dsa.MarshalPrivateKey(test.want.Sparkle.DSAKey)
				b = pem.EncodeToMemory(&pem.Block{Type: "PRIVATE KEY", Bytes: b})
				os.WriteFile(filepath.Join(dir, "dsa_key"), b, 0o600)
			}
			if test.want.Sparkle.Ed25519Key != nil {
				b, _ := ed25519.MarshalPrivateKey(test.want.Sparkle.Ed25519Key)
				b = pem.EncodeToMemory(&pem.Block{Type: "PRIVATE KEY", Bytes: b})
				os.WriteFile(filepath.Join(dir, "ed25519_key"), b, 0o600)
			}
		}

		path := filepath.Join(t.TempDir(), "appcast.yml")
		os.WriteFile(path, []byte(test.in), os.ModePerm)

		c, err := pipe.Load(path)
		if err != nil {
			t.Error(err)
			continue
		}

		if diff := cmp.Diff(test.want, c, cmp.Exporter(func(t reflect.Type) bool { return true })); diff != "" {
			t.Error(i, diff)
		}
	}
}
