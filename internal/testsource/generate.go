package testsource

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/goreleaser/nfpm/v2"
	"github.com/goreleaser/nfpm/v2/files"

	_ "github.com/goreleaser/nfpm/v2/apk"  // apk packager
	_ "github.com/goreleaser/nfpm/v2/arch" // archlinux packager
	_ "github.com/goreleaser/nfpm/v2/deb"  // deb packager
	_ "github.com/goreleaser/nfpm/v2/rpm"  // rpm packager
)

type Generator func(dir, version string) error

func getConfig(version, arch string) nfpm.Config {
	return nfpm.Config{
		Info: nfpm.Info{
			Name:        "appcast-test",
			Arch:        arch,
			Platform:    "linux",
			Version:     version,
			Section:     "utils",
			Priority:    "optional",
			Maintainer:  "Test User <test@example.com>",
			Description: "This is a test.\nIt does nothing.\n\nAbsolutely nothing.",
			Vendor:      "Test Company",
			Homepage:    "http://example.com",
			License:     "MIT",
			Overridables: nfpm.Overridables{
				Replaces:   []string{"appcast-test-old"},
				Provides:   []string{"appcast-test-alt"},
				Depends:    []string{"bash"},
				Recommends: []string{"git"},
				Suggests:   []string{"wget"},
				Conflicts:  []string{"appcast-test-new"},
				Contents: files.Contents{
					{
						Source:      "./appcast-test",
						Destination: "/usr/bin/appcast-test",
					},
				},
			},
		},
	}
}

func GenerateLinux(packager string) Generator {
	return func(dir, version string) error {
		for _, arch := range []string{"amd64", "386"} {
			config := getConfig(version, arch)

			srcDir, err := os.MkdirTemp("", "")
			if err != nil {
				return err
			}
			defer os.RemoveAll(srcDir)

			binPath := filepath.Join(srcDir, "appcast-test")
			bin := []byte(fmt.Sprintf(`#/bin/bash\n\necho "test %s"\n`, config.Version))
			if err = os.WriteFile(filepath.Join(srcDir, "appcast-test"), bin, 0o755); err != nil { //nolint:gosec
				return err
			}
			config.Contents[0].Source = binPath

			info, err := config.Get(packager)
			if err != nil {
				return err
			}

			info = nfpm.WithDefaults(info)

			pkg, err := nfpm.Get(packager)
			if err != nil {
				return err
			}

			config.Target = filepath.Join(dir, pkg.ConventionalFileName(info))

			if err = os.MkdirAll(filepath.Dir(config.Target), 0o755); err != nil {
				return err
			}
			f, err := os.Create(config.Target)
			if err != nil {
				return err
			}
			defer f.Close()

			if err := pkg.Package(info, f); err != nil {
				return err
			}
		}

		return nil
	}
}
