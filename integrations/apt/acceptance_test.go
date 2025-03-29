//go:build acceptance

package apt_test

import (
	"testing"

	"github.com/kubri/kubri/integrations/apt"
	"github.com/kubri/kubri/internal/emulator"
	"github.com/kubri/kubri/pkg/crypto/pgp"
	source "github.com/kubri/kubri/source/file"
	ftarget "github.com/kubri/kubri/target/file"
)

func TestAcceptance(t *testing.T) {
	distros := []struct {
		name  string
		image string
	}{
		{"Debian", "debian:latest"},
		{"Ubuntu", "ubuntu:latest"},
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
			t.Parallel()

			dir := t.TempDir()
			pgpKey, _ := pgp.NewPrivateKey("test", "test@example.com")
			key, _ := pgp.MarshalPublicKey(pgp.Public(pgpKey))
			src, _ := source.New(source.Config{Path: "../../testdata"})
			tgt, _ := ftarget.New(ftarget.Config{Path: dir})
			url := emulator.FileServer(t, dir)
			c := emulator.Image(t, distro.image, emulator.WithEnv("DEBIAN_FRONTEND", "noninteractive"))
			c.Exec(t, "apt-get update -q && apt-get install -yq --no-install-recommends gpg")
			c.CopyToContainer(t.Context(), key, "kubri-test.asc", 0o644)
			c.Exec(t, "gpg --dearmor --yes --output /usr/share/keyrings/kubri-test.gpg < kubri-test.asc")
			c.Exec(t, "echo 'deb [signed-by=/usr/share/keyrings/kubri-test.gpg] "+url+" stable main' > /etc/apt/sources.list.d/kubri-test.list")

			for i, test := range tests {
				t.Run(test.name, func(t *testing.T) {
					config := &apt.Config{Source: src, Target: tgt, Version: test.version, PGPKey: pgpKey}
					if err := apt.Build(t.Context(), config); err != nil {
						t.Fatal(err)
					}

					c.Exec(t, "apt-get update -q")

					if i == 0 {
						c.Exec(t, "apt-get install -yq --no-install-recommends kubri-test")
					} else {
						c.Exec(t, "apt-get upgrade -yq")
					}

					if v := c.Exec(t, "kubri-test"); v != test.version {
						t.Fatalf("expected version %q got %q", test, v)
					}
				})

				if t.Failed() {
					t.FailNow()
				}
			}
		})
	}
}
