package appinstaller_test

import (
	"io"
	"testing"

	"github.com/google/go-cmp/cmp"

	"github.com/kubri/kubri/integrations/appinstaller"
	source "github.com/kubri/kubri/source/file"
	target "github.com/kubri/kubri/target/file"
)

func TestBuild(t *testing.T) {
	src, _ := source.New(source.Config{Path: "../../testdata", URL: "https://dl.example.com/"})
	tgt, _ := target.New(target.Config{Path: t.TempDir(), URL: "https://example.com"})

	tests := []struct {
		name   string
		config *appinstaller.Config
		want   string
	}{
		{
			name: "Test-x64.appinstaller",
			config: &appinstaller.Config{
				Source: src,
				Target: tgt,
			},
			want: `<?xml version="1.0" encoding="UTF-8"?>
<AppInstaller xmlns="http://schemas.microsoft.com/appx/appinstaller/2017" Version="1.0.0.1" Uri="https://example.com/Test-x64.appinstaller">
	<MainPackage Name="Test" Publisher="CN=Test" Version="1.0.0.1" ProcessorArchitecture="x64" Uri="https://dl.example.com/v1.0.0/test.msix" />
</AppInstaller>`,
		},
		{
			name: "Test.appinstaller",
			config: &appinstaller.Config{
				Source: src,
				Target: tgt,
				OnLaunch: &appinstaller.OnLaunchConfig{
					ShowPrompt: true,
				},
			},
			want: `<?xml version="1.0" encoding="UTF-8"?>
<AppInstaller xmlns="http://schemas.microsoft.com/appx/appinstaller/2018" Version="1.0.0.1" Uri="https://example.com/Test.appinstaller">
	<MainBundle Name="Test" Publisher="CN=Test" Version="1.0.0.1" Uri="https://dl.example.com/v1.0.0/test.msixbundle" />
	<UpdateSettings>
		<OnLaunch ShowPrompt="true" />
	</UpdateSettings>
</AppInstaller>`,
		},
		{
			name: "Test-x64.appinstaller",
			config: &appinstaller.Config{
				Source:                  src,
				Target:                  tgt,
				AutomaticBackgroundTask: true,
			},
			want: `<?xml version="1.0" encoding="UTF-8"?>
<AppInstaller xmlns="http://schemas.microsoft.com/appx/appinstaller/2017/2" Version="1.0.0.1" Uri="https://example.com/Test-x64.appinstaller">
	<MainPackage Name="Test" Publisher="CN=Test" Version="1.0.0.1" ProcessorArchitecture="x64" Uri="https://dl.example.com/v1.0.0/test.msix" />
	<UpdateSettings>
		<AutomaticBackgroundTask />
	</UpdateSettings>
</AppInstaller>`,
		},
		{
			name: "Test-x64.appinstaller",
			config: &appinstaller.Config{
				Source: src,
				Target: tgt,
				OnLaunch: &appinstaller.OnLaunchConfig{
					HoursBetweenUpdateChecks: 12,
					UpdateBlocksActivation:   true,
					ShowPrompt:               true,
				},
				AutomaticBackgroundTask:   true,
				ForceUpdateFromAnyVersion: true,
			},
			want: `<?xml version="1.0" encoding="UTF-8"?>
<AppInstaller xmlns="http://schemas.microsoft.com/appx/appinstaller/2018" Version="1.0.0.1" Uri="https://example.com/Test-x64.appinstaller">
	<MainPackage Name="Test" Publisher="CN=Test" Version="1.0.0.1" ProcessorArchitecture="x64" Uri="https://dl.example.com/v1.0.0/test.msix" />
	<UpdateSettings>
		<OnLaunch HoursBetweenUpdateChecks="12" ShowPrompt="true" UpdateBlocksActivation="true" />
		<AutomaticBackgroundTask />
		<ForceUpdateFromAnyVersion>true</ForceUpdateFromAnyVersion>
	</UpdateSettings>
</AppInstaller>`,
		},
		{
			name: "Test.appinstaller",
			config: &appinstaller.Config{
				Source:         src,
				Target:         tgt,
				UploadPackages: true,
			},
			want: `<?xml version="1.0" encoding="UTF-8"?>
<AppInstaller xmlns="http://schemas.microsoft.com/appx/appinstaller/2017" Version="1.0.0.1" Uri="https://example.com/Test.appinstaller">
	<MainBundle Name="Test" Publisher="CN=Test" Version="1.0.0.1" Uri="https://example.com/test.msixbundle" />
</AppInstaller>`,
		},
	}

	for _, test := range tests {
		err := appinstaller.Build(t.Context(), test.config)
		if err != nil {
			t.Fatal(err)
		}

		r, err := test.config.Target.NewReader(t.Context(), test.name)
		if err != nil {
			t.Fatal(err)
		}
		defer r.Close()

		got, err := io.ReadAll(r)
		if err != nil {
			t.Fatal(err)
		}

		if diff := cmp.Diff(test.want, string(got)); diff != "" {
			t.Fatal(diff)
		}

		if test.config.UploadPackages {
			for _, ext := range []string{".msixbundle", ".msix"} {
				if _, err = tgt.NewReader(t.Context(), "test"+ext); err != nil {
					t.Fatal(err)
				}
			}
		}
	}
}
