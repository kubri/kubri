//go:build acceptance

package apk_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/abemedia/appcast/integrations/apk"
	"github.com/abemedia/appcast/internal/emulator"
	"github.com/abemedia/appcast/pkg/crypto/rsa"
	source "github.com/abemedia/appcast/source/file"
	target "github.com/abemedia/appcast/target/file"
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
			s := httptest.NewServer(http.FileServer(http.Dir(dir)))
			c := emulator.Image(t, distro.image)

			for i, test := range tests {
				t.Run(test.name, func(t *testing.T) {
					config := &apk.Config{Source: src, Target: tgt, Version: test.version, RSAKey: rsaKey, KeyName: "test@example.com.rsa.pub"}
					if err := apk.Build(context.Background(), config); err != nil {
						t.Fatal(err)
					}

					if i == 0 {
						c.Exec(t, "echo '"+s.URL+"' >> /etc/apk/repositories")
						c.Exec(t, "apk add --no-cache wget")
						c.Exec(t, "wget -q -O /etc/apk/keys/"+config.KeyName+" "+s.URL+"/"+config.KeyName)
						c.Exec(t, "apk add --no-cache appcast-test")
					} else {
						c.Exec(t, "apk upgrade --no-cache appcast-test")
					}

					if v := c.Exec(t, "appcast-test"); v != test.version {
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
