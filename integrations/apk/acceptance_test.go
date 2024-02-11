//go:build acceptance

package apk_test

import (
	"context"
	"testing"

	"github.com/kubri/kubri/integrations/apk"
	"github.com/kubri/kubri/internal/emulator"
	"github.com/kubri/kubri/pkg/crypto/rsa"
	source "github.com/kubri/kubri/source/file"
	target "github.com/kubri/kubri/target/file"
)

func TestAcceptance(t *testing.T) {
	distros := []struct {
		name  string
		image string
	}{
		{"Alpine 3", "alpine:3"},
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
			rsaKey, _ := rsa.NewPrivateKey()
			src, _ := source.New(source.Config{Path: "../../testdata"})
			tgt, _ := target.New(target.Config{Path: dir})
			url := emulator.FileServer(t, dir)
			c := emulator.Image(t, distro.image)

			for i, test := range tests {
				t.Run(test.name, func(t *testing.T) {
					config := &apk.Config{Source: src, Target: tgt, Version: test.version, RSAKey: rsaKey, KeyName: "test@example.com"}
					if err := apk.Build(context.Background(), config); err != nil {
						t.Fatal(err)
					}

					if i == 0 {
						c.Exec(t, "echo '"+url+"' >> /etc/apk/repositories")
						c.Exec(t, "apk add --no-cache wget")
						c.Exec(t, "wget -q -O /etc/apk/keys/"+config.KeyName+".rsa.pub "+url+"/"+config.KeyName+".rsa.pub")
						c.Exec(t, "apk add --no-cache kubri-test")
					} else {
						c.Exec(t, "apk upgrade --no-cache kubri-test")
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
