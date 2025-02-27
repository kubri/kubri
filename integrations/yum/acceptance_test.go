//go:build acceptance

package yum_test

import (
	"fmt"
	"testing"

	"github.com/kubri/kubri/integrations/yum"
	"github.com/kubri/kubri/internal/emulator"
	"github.com/kubri/kubri/pkg/crypto/pgp"
	source "github.com/kubri/kubri/source/file"
	target "github.com/kubri/kubri/target/file"
)

const conf = `[kubri-test]
name=kubri-test
baseurl=%s
enabled=1
gpgcheck=0
repo_gpgcheck=1
gpgkey=%s/repodata/repomd.xml.key`

func TestAcceptance(t *testing.T) {
	distros := []struct {
		name  string
		image string
		pkg   string
	}{
		{"RHEL 9", "registry.access.redhat.com/ubi9/ubi:latest", "dnf"},
		{"RHEL 8", "registry.access.redhat.com/ubi8/ubi:latest", "dnf"},
		{"Fedora 39", "fedora:39", "dnf"},
		{"Fedora 38", "fedora:38", "dnf"},
		{"openSUSE Leap 15", "opensuse/leap:15", "zypper"},
	}

	tests := []struct {
		name    string
		version string
	}{
		{"Install", "v1.0.0"},
		{"Update", "v2.0.0"},
	}

	for _, distro := range distros {
		t.Run(distro.name, func(t *testing.T) {
			dir := t.TempDir()
			pgpKey, _ := pgp.NewPrivateKey("test", "test@example.com")
			src, _ := source.New(source.Config{Path: "../../testdata"})
			tgt, _ := target.New(target.Config{Path: dir})
			url := emulator.FileServer(t, dir)
			c := emulator.Image(t, distro.image)

			for i, test := range tests {
				t.Run(test.name, func(t *testing.T) {
					config := &yum.Config{Source: src, Target: tgt, Version: test.version, PGPKey: pgpKey}
					if err := yum.Build(t.Context(), config); err != nil {
						t.Fatal(err)
					}

					switch distro.pkg {
					case "dnf":
						if i == 0 {
							c.Exec(t, "echo '"+fmt.Sprintf(conf, url, url)+"' > /etc/yum.repos.d/kubri-test.repo")
							c.Exec(t, "dnf install -yq kubri-test")
						} else {
							c.Exec(t, "dnf clean expire-cache")
							c.Exec(t, "dnf update -yq kubri-test")
						}
					case "zypper":
						if i == 0 {
							c.Exec(t, "zypper addrepo --refresh "+url+" kubri-test")
							c.Exec(t, "zypper --gpg-auto-import-keys refresh")
							c.Exec(t, "zypper --non-interactive install kubri-test")
						} else {
							c.Exec(t, "zypper refresh")
							c.Exec(t, "zypper --non-interactive update kubri-test")
						}
					}

					if v := c.Exec(t, "kubri-test"); v != test.version {
						t.Fatalf("expected version %q got %q", test.version, v)
					}
				})

				if t.Failed() {
					t.FailNow()
				}
			}
		})
	}
}
