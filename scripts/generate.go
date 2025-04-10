package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/goreleaser/nfpm/v2"
	"github.com/goreleaser/nfpm/v2/files"

	_ "github.com/goreleaser/nfpm/v2/apk"  // apk packager
	_ "github.com/goreleaser/nfpm/v2/arch" // archlinux packager
	_ "github.com/goreleaser/nfpm/v2/deb"  // deb packager
	_ "github.com/goreleaser/nfpm/v2/ipk"  // ipk packager
	_ "github.com/goreleaser/nfpm/v2/rpm"  // rpm packager
)

//nolint:gochecknoglobals
var (
	versions = []string{"v1.0.0", "v1.1.0-beta", "v1.1.0", "v2.0.0"}
	archs    = []string{"amd64", "386"}
	formats  = []string{"deb", "rpm", "apk", "archlinux"}

	config = nfpm.Config{
		Info: nfpm.Info{
			MTime:       time.Date(2023, 11, 19, 23, 37, 12, 0, time.UTC),
			Name:        "kubri-test",
			Platform:    "linux",
			Section:     "utils",
			Priority:    "optional",
			Maintainer:  "Test User <test@example.com>",
			Description: "This is a test.\nIt does nothing.\n\nAbsolutely nothing.",
			Vendor:      "Test Company",
			Homepage:    "http://example.com",
			License:     "MIT",
			Overridables: nfpm.Overridables{
				Replaces:   []string{"kubri-test-old"},
				Provides:   []string{"kubri-test-alt"},
				Depends:    []string{"bash"},
				Recommends: []string{"git"},
				Suggests:   []string{"wget"},
				Conflicts:  []string{"kubri-test-new"},
				Contents: files.Contents{
					{
						Source:      "./kubri-test",
						Destination: "/usr/bin/kubri-test",
					},
				},
				Deb: nfpm.Deb{Compression: "xz"},
				RPM: nfpm.RPM{Group: "default"},
				ArchLinux: nfpm.ArchLinux{
					Packager: "Kubri <info@kubri.dev>",
				},
			},
		},
	}
)

func main() {
	c := config
	for _, version := range versions {
		c.Version = version
		for _, format := range formats {
			for _, arch := range archs {
				c.Arch = arch
				if err := buildPackages(format, c); err != nil {
					fmt.Printf("failed to build %s package: %s\n", formats, err) //nolint:forbidigo
					os.Exit(1)
				}
			}
		}
	}
}

func buildPackages(packager string, config nfpm.Config) error {
	srcDir, err := os.MkdirTemp("", "")
	if err != nil {
		return err
	}
	defer os.RemoveAll(srcDir)

	binPath := filepath.Join(srcDir, "kubri-test")
	bin := fmt.Appendf(nil, "#/bin/bash\n\necho %q\n", config.Version)
	if err = os.WriteFile(filepath.Join(srcDir, "kubri-test"), bin, 0o755); err != nil { //nolint:gosec
		return err
	}
	config.Contents[0].Source = binPath

	info, err := config.Get(packager)
	if err != nil {
		return err
	}

	if packager == "apk" {
		d := strings.Split(info.Description, "\n")
		for i, l := range d {
			d[i] = strings.TrimSpace(l)
		}
		info.Description = strings.Join(d, " ")
	}

	info = nfpm.WithDefaults(info)

	pkg, err := nfpm.Get(packager)
	if err != nil {
		return err
	}

	config.Target = filepath.Join("testdata", config.Version, pkg.ConventionalFileName(info))

	if err = os.MkdirAll(filepath.Dir(config.Target), 0o750); err != nil {
		return err
	}

	f, err := os.Create(config.Target)
	if err != nil {
		return err
	}
	defer f.Close()

	if err := pkg.Package(info, f); err != nil {
		_ = os.Remove(config.Target)
		return err
	}

	fmt.Printf("created package: %s\n", config.Target) //nolint:forbidigo
	return f.Close()
}
