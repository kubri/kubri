//go:build acceptance

package arch_test

import (
	"fmt"
	"testing"

	"github.com/kubri/kubri/integrations/arch"
	"github.com/kubri/kubri/internal/emulator"
	"github.com/kubri/kubri/pkg/crypto/pgp"
	source "github.com/kubri/kubri/source/file"
	target "github.com/kubri/kubri/target/file"
)

const conf = `[kubri-test]
SigLevel = Required
Server = %s/$arch
`

func TestAcceptance(t *testing.T) {
	distros := []struct {
		name  string
		image string
	}{
		{"Arch Linux", "archlinux:latest"},
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
			tgt, _ := target.New(target.Config{Path: dir})
			url := emulator.FileServer(t, dir)
			c := emulator.Image(t, distro.image)

			c.CopyToContainer(t.Context(), key, "kubri-test.asc", 0o644)
			c.Exec(t, "pacman-key --init")
			c.Exec(t, "pacman-key --populate")
			c.Exec(t, "pacman-key --add kubri-test.asc")
			c.Exec(t, `pacman-key --lsign-key "$(gpg --with-colons --import-options show-only --import kubri-test.asc | awk -F: '/^fpr:/ {print $10; exit}')"`)

			for _, test := range tests {
				t.Run(test.name, func(t *testing.T) {
					config := &arch.Config{
						RepoName: "kubri-test",
						Source:   src,
						Version:  test.version,
						Target:   tgt,
						PGPKey:   pgpKey,
					}
					if err := arch.Build(t.Context(), config); err != nil {
						t.Fatal(err)
					}

					c.Exec(t, "echo '"+fmt.Sprintf(conf, url)+"' >> /etc/pacman.conf")
					c.Exec(t, "pacman --noconfirm -Syu kubri-test")

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
