//go:build acceptance

package apt_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/abemedia/appcast/integrations/apt"
	"github.com/abemedia/appcast/internal/emulator"
	"github.com/abemedia/appcast/pkg/crypto/pgp"
	source "github.com/abemedia/appcast/source/file"
	ftarget "github.com/abemedia/appcast/target/file"
)

func TestAcceptance(t *testing.T) {
	distros := []struct {
		name  string
		image string
	}{
		{"Debian 12", "debian:12"},
		{"Debian 11", "debian:11"},
		{"Ubuntu 22.04", "ubuntu:jammy"},
		{"Ubuntu 20.04", "ubuntu:focal"},
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
			key, _ := pgp.MarshalPublicKey(pgp.Public(pgpKey))
			src, _ := source.New(source.Config{Path: "../../testdata"})
			tgt, _ := ftarget.New(ftarget.Config{Path: dir})
			s := httptest.NewServer(http.FileServer(http.Dir(dir)))
			c := emulator.Build(t, `
				FROM `+distro.image+`
				ENV DEBIAN_FRONTEND=noninteractive
				RUN apt-get update && apt-get install -y --no-install-recommends gpg apt-utils
				ENTRYPOINT ["tail", "-f", "/dev/null"]
			`)

			c.CopyToContainer(context.Background(), key, "appcast-test.asc", 0o644)
			c.Exec(t, "gpg --dearmor --yes --output /usr/share/keyrings/appcast-test.gpg < appcast-test.asc")
			c.Exec(t, "echo 'deb [signed-by=/usr/share/keyrings/appcast-test.gpg] "+s.URL+" stable main' > /etc/apt/sources.list.d/appcast-test.list")

			for i, test := range tests {
				t.Run(test.name, func(t *testing.T) {
					config := &apt.Config{Source: src, Target: tgt, Version: test.version, PGPKey: pgpKey}
					if err := apt.Build(context.Background(), config); err != nil {
						t.Fatal(err)
					}

					c.Exec(t, "apt-get update -qq")

					if i == 0 {
						c.Exec(t, "apt-get install -qq -y --no-install-recommends appcast-test")
					} else {
						c.Exec(t, "apt-get upgrade -qq -y")
					}

					if v := c.Exec(t, "appcast-test"); v != test.version {
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
