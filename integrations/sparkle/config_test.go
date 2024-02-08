package sparkle_test

import (
	"testing"

	"github.com/google/go-cmp/cmp"

	"github.com/kubri/kubri/integrations/sparkle"
)

func TestConfigGetOptions(t *testing.T) {
	c := &sparkle.Config{
		Settings: []sparkle.Rule{
			{
				OS: sparkle.MacOS,
				Settings: &sparkle.Settings{
					MinimumSystemVersion: "10.13.0",
				},
			},
			{
				OS: sparkle.Windows64,
				Settings: &sparkle.Settings{
					InstallerArguments: "/passive",
				},
			},
			{
				Version: "> 1.0.0",
				Settings: &sparkle.Settings{
					MinimumAutoupdateVersion: "1.0.0",
				},
			},
			{
				Version: ">= 1.1.0, < 1.2.0",
				Settings: &sparkle.Settings{
					CriticalUpdate: true,
				},
			},
			{
				Version: "v1.2.1",
				Settings: &sparkle.Settings{
					CriticalUpdateBelowVersion: "1.1.0",
					MinimumAutoupdateVersion:   "1.1.0",
				},
			},
		},
	}

	tests := []struct {
		version string
		os      sparkle.OS
		want    *sparkle.Settings
	}{
		{
			version: "1.0.0",
			os:      sparkle.Windows64,
			want: &sparkle.Settings{
				InstallerArguments: "/passive",
			},
		},
		{
			version: "1.0.0",
			os:      sparkle.MacOS,
			want: &sparkle.Settings{
				MinimumSystemVersion: "10.13.0",
			},
		},
		{
			version: "1.1.1",
			os:      sparkle.Windows64,
			want: &sparkle.Settings{
				InstallerArguments:       "/passive",
				MinimumAutoupdateVersion: "1.0.0",
				CriticalUpdate:           true,
			},
		},
		{
			version: "1.2.1",
			os:      sparkle.Windows,
			want: &sparkle.Settings{
				MinimumAutoupdateVersion:   "1.1.0",
				CriticalUpdateBelowVersion: "1.1.0",
			},
		},
	}

	for _, test := range tests {
		opt, err := sparkle.GetSettings(c.Settings, test.version, test.os)
		if err != nil {
			t.Error(test.version, test.os, err)
			continue
		}

		if diff := cmp.Diff(test.want, opt); diff != "" {
			t.Error(test.version, test.os, diff)
		}
	}
}
