package pipe_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/abemedia/appcast/integrations/appinstaller"
	"github.com/abemedia/appcast/integrations/apt"
	"github.com/abemedia/appcast/integrations/sparkle"
	"github.com/abemedia/appcast/internal/test"
	"github.com/abemedia/appcast/pkg/crypto/dsa"
	"github.com/abemedia/appcast/pkg/crypto/ed25519"
	"github.com/abemedia/appcast/pkg/pipe"
	source "github.com/abemedia/appcast/source/file"
	target "github.com/abemedia/appcast/target/file"
	"github.com/google/go-cmp/cmp"
)

func TestPipe(t *testing.T) {
	dir := t.TempDir()

	src, err := source.New(source.Config{Path: dir})
	if err != nil {
		t.Fatal(err)
	}

	tgt, err := target.New(target.Config{Path: dir})
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
appinstaller: {}
apt: {}
sparkle: {}
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
					FileName:    "appcast.xml",
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
  hours_between_update_checks: 12
  update_blocks_activation: true
  show_prompt: true
  automatic_background_task: true
  force_update_from_any_version: true	
apt:
  folder: .
sparkle:
  folder: .
  filename: updates.xml
  title: foo
  description: bar
  params:
    - os: windows
      installer_arguments: /passive
    - os: macos
      minimum_system_version: '10.13.0'
    - version: '1.0.0'
      critical_update: true
    - version: '> 1.0.0'
      critical_update_below_version: '1.0.0'
      minimum_autoupdate_version: '1.0.0'
    - version: '1.1.0'
      ignore_skipped_upgrades_below_version: '1.1.0'
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

	opts := test.ExportAll()

	for i, test := range tests {
		dir := t.TempDir()
		t.Setenv("APPCAST_PATH", dir)

		if test.want.Sparkle != nil {
			if test.want.Sparkle.DSAKey != nil {
				b, _ := dsa.MarshalPrivateKey(test.want.Sparkle.DSAKey)
				os.WriteFile(filepath.Join(dir, "dsa_key"), b, 0o600)
			}
			if test.want.Sparkle.Ed25519Key != nil {
				b, _ := ed25519.MarshalPrivateKey(test.want.Sparkle.Ed25519Key)
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

		if diff := cmp.Diff(test.want, c, opts); diff != "" {
			t.Error(i, diff)
		}
	}
}
