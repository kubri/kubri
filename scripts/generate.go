package main

import (
	"context"
	"errors"
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

	"github.com/kubri/kubri/integrations/apk"
	"github.com/kubri/kubri/integrations/apt"
	"github.com/kubri/kubri/integrations/yum"
	source "github.com/kubri/kubri/source/file"
	target "github.com/kubri/kubri/target/file"
)

//nolint:gochecknoglobals
var (
	versions = []string{"v1.0.0", "v1.1.0-beta", "v1.1.0", "v2.0.0"}
	archs    = []string{"amd64", "386"}
	formats  = []string{"deb", "rpm", "apk"}

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

	err := errors.Join(apkGolden(), aptGolden(), yumGolden())
	if err != nil {
		fmt.Printf("failed to generate golden: %s\n", err) //nolint:forbidigo
		os.Exit(1)
	}
}

func buildPackages(packager string, config nfpm.Config) error {
	srcDir, err := os.MkdirTemp("", "")
	if err != nil {
		return err
	}
	defer os.RemoveAll(srcDir)

	binPath := filepath.Join(srcDir, "kubri-test")
	bin := []byte(fmt.Sprintf("#/bin/bash\n\necho %q\n", config.Version))
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
		os.Remove(config.Target)
		return err
	}

	fmt.Printf("created package: %s\n", config.Target) //nolint:forbidigo
	return f.Close()
}

func aptGolden() error {
	path := filepath.Join("integrations", "apt", "testdata")

	if err := os.RemoveAll(path); err != nil {
		return err
	}

	src, err := source.New(source.Config{Path: "testdata"})
	if err != nil {
		return err
	}

	tgt, err := target.New(target.Config{Path: path})
	if err != nil {
		return err
	}

	err = apt.Build(context.Background(), &apt.Config{Source: src, Target: tgt})
	if err != nil {
		return err
	}

	err = os.RemoveAll(filepath.Join(path, "pool"))
	if err != nil {
		return err
	}

	return filepath.WalkDir(path, func(path string, d os.DirEntry, err error) error {
		if err != nil || d.IsDir() {
			return err
		}
		if filepath.Ext(path) != "" {
			return os.Remove(path)
		}
		return nil
	})
}

func yumGolden() error {
	path := filepath.Join("integrations", "yum", "testdata")

	if err := os.RemoveAll(path); err != nil {
		return err
	}

	src, err := source.New(source.Config{Path: "testdata"})
	if err != nil {
		return err
	}

	tgt, err := target.New(target.Config{Path: path})
	if err != nil {
		return err
	}

	err = yum.Build(context.Background(), &yum.Config{Source: src, Target: tgt})
	if err != nil {
		return err
	}

	return os.RemoveAll(filepath.Join(path, "Packages"))
}

func apkGolden() error {
	dir := filepath.Join("integrations", "apk", "testdata")

	if err := os.RemoveAll(dir); err != nil {
		return err
	}

	src, err := source.New(source.Config{Path: "testdata"})
	if err != nil {
		return err
	}

	tgt, err := target.New(target.Config{Path: dir})
	if err != nil {
		return err
	}

	err = apk.Build(context.Background(), &apk.Config{Source: src, Target: tgt})
	if err != nil {
		return err
	}

	return filepath.WalkDir(dir, func(path string, d os.DirEntry, err error) error {
		if err != nil || d.IsDir() {
			return err
		}
		if filepath.Ext(path) == ".apk" {
			return os.Remove(path)
		}
		return nil
	})
}
