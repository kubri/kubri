package config_test

import (
	"testing"

	"github.com/kubri/kubri/integrations/appinstaller"
	"github.com/kubri/kubri/pkg/config"
	source "github.com/kubri/kubri/source/file"
	target "github.com/kubri/kubri/target/file"
)

func TestAppInstaller(t *testing.T) {
	dir := t.TempDir()
	src, _ := source.New(source.Config{Path: dir})
	tgt, _ := target.New(target.Config{Path: dir})

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
				appinstaller:
					disabled: true
			`,
			want: &config.Config{},
		},
		{
			desc: "defaults",
			in: `
				source:
					type: file
					path: ` + dir + `
				target:
					type: file
					path: ` + dir + `
				appinstaller: {}
			`,
			want: &config.Config{
				Appinstaller: &appinstaller.Config{
					Source: src,
					Target: tgt.Sub("appinstaller"),
				},
			},
		},
		{
			desc: "full",
			in: `
				version: latest
				prerelease: true
				upload-packages: true
				source:
					type: file
					path: ` + dir + `
				target:
					type: file
					path: ` + dir + `
				appinstaller:
					folder: test
					on-launch:
						hours-between-update-checks: 12
						show-prompt: true
						update-blocks-activation: true
					automatic-background-task: true
					force-update-from-any-version: true
			`,
			want: &config.Config{
				Appinstaller: &appinstaller.Config{
					OnLaunch: &appinstaller.OnLaunchConfig{
						HoursBetweenUpdateChecks: 12,
						ShowPrompt:               true,
						UpdateBlocksActivation:   true,
					},
					AutomaticBackgroundTask:   true,
					ForceUpdateFromAnyVersion: true,

					Source:         src,
					Target:         tgt.Sub("test"),
					Version:        "latest",
					Prerelease:     true,
					UploadPackages: true,
				},
			},
		},
		{
			desc: "invalid folder",
			in: `
				version: latest
				prerelease: true
				source:
					type: file
					path: ` + dir + `
				target:
					type: file
					path: ` + dir + `
				appinstaller:
					folder: '*'
			`,
			err: &config.Error{Errors: []string{"appinstaller.folder must be a valid folder name"}},
		},
		{
			desc: "hours-between-update-checks below 0",
			in: `
				source:
					type: file
					path: ` + dir + `
				target:
					type: file
					path: ` + dir + `
				appinstaller:
					on-launch:
						hours-between-update-checks: -1
			`,
			err: &config.Error{Errors: []string{"appinstaller.on-launch.hours-between-update-checks must be 0 or greater"}},
		},
		{
			desc: "hours-between-update-checks above 255",
			in: `
				source:
					type: file
					path: ` + dir + `
				target:
					type: file
					path: ` + dir + `
				appinstaller:
					on-launch:
						hours-between-update-checks: 256
			`,
			err: &config.Error{Errors: []string{"appinstaller.on-launch.hours-between-update-checks must be 255 or less"}},
		},
	})
}
