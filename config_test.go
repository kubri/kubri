package appcast_test

import (
	"testing"

	"github.com/abemedia/appcast"
	"github.com/google/go-cmp/cmp"
)

func TestConfigGetOptions(t *testing.T) {
	c := &appcast.Config{
		Settings: []appcast.Rule{
			{
				OS: appcast.MacOS,
				Settings: &appcast.Settings{
					MinimumSystemVersion: "10.13.0",
				},
			},
			{
				OS: appcast.Windows64,
				Settings: &appcast.Settings{
					InstallerArguments: "/passive",
				},
			},
			{
				Version: "> 1.0.0",
				Settings: &appcast.Settings{
					MinimumAutoupdateVersion: "1.0.0",
				},
			},
			{
				Version: ">= 1.1.0, < 1.2.0",
				Settings: &appcast.Settings{
					CriticalUpdate: true,
				},
			},
			{
				Version: "~> 1.2.0",
				Settings: &appcast.Settings{
					CriticalUpdateBelowVersion: "1.1.0",
				},
			},
		},
	}

	tests := []struct {
		version string
		os      appcast.OS
		want    *appcast.Settings
	}{
		{
			version: "1.0.0",
			os:      appcast.Windows64,
			want: &appcast.Settings{
				InstallerArguments: "/passive",
			},
		},
		{
			version: "1.0.0",
			os:      appcast.MacOS,
			want: &appcast.Settings{
				MinimumSystemVersion: "10.13.0",
			},
		},
		{
			version: "1.1.1",
			os:      appcast.Windows64,
			want: &appcast.Settings{
				InstallerArguments:       "/passive",
				MinimumAutoupdateVersion: "1.0.0",
				CriticalUpdate:           true,
			},
		},
		{
			version: "1.2.1",
			os:      appcast.Windows,
			want: &appcast.Settings{
				MinimumAutoupdateVersion:   "1.0.0",
				CriticalUpdateBelowVersion: "1.1.0",
			},
		},
	}

	for _, test := range tests {
		opt, err := appcast.GetSettings(c.Settings, test.version, test.os)
		if err != nil {
			t.Error(test.version, test.os, err)
			continue
		}

		if diff := cmp.Diff(test.want, opt); diff != "" {
			t.Error(test.version, test.os, diff)
		}
	}
}
