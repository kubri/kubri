package config_test

import (
	"encoding/pem"
	"os"
	"path"
	"reflect"
	"testing"

	"github.com/abemedia/appcast/integrations/sparkle"
	"github.com/abemedia/appcast/pkg/config"
	"github.com/abemedia/appcast/pkg/crypto/dsa"
	"github.com/abemedia/appcast/pkg/crypto/ed25519"
	"github.com/abemedia/appcast/source"
	_ "github.com/abemedia/appcast/source/blob/memory"
	"github.com/abemedia/appcast/target"
	_ "github.com/abemedia/appcast/target/blob/memory"
	"github.com/google/go-cmp/cmp"
)

func TestSparkle(t *testing.T) {
	t.Setenv("APPCAST_PATH", t.TempDir())

	src, err := source.Open("mem://")
	if err != nil {
		t.Fatal(err)
	}

	tgt, err := target.Open("mem://")
	if err != nil {
		t.Fatal(err)
	}

	dsaKey, _ := dsa.NewPrivateKey()
	edKey, _ := ed25519.NewPrivateKey()

	tests := []struct {
		in   *config.Config
		want *sparkle.Config
	}{
		{
			in: &config.Config{
				Title:       "test",
				Description: "test",
				Source:      config.Source{Repo: "mem://"},
				Target:      config.Target{Repo: "mem://"},
			},
			want: &sparkle.Config{
				Title:       "test",
				Description: "test",
				Source:      src,
				Target:      tgt.Sub("sparkle"),
				FileName:    "sparkle.xml",
			},
		},
		{
			in: &config.Config{
				Title:       "test",
				Description: "test",
				Source:      config.Source{Repo: "mem://", Version: "latest"},
				Target:      config.Target{Repo: "mem://", Flat: true},
				Sparkle: config.Sparkle{
					Title:       "foo",
					Description: "bar",
					FileName:    "updates.xml",
				},
			},
			want: &sparkle.Config{
				Title:       "foo",
				Description: "bar",
				Source:      src,
				Target:      tgt,
				FileName:    "updates.xml",
				DSAKey:      dsaKey,
				Ed25519Key:  edKey,
				Version:     "latest",
			},
		},
	}

	for _, test := range tests {
		if test.want.DSAKey != nil {
			b, _ := dsa.MarshalPrivateKey(test.want.DSAKey)
			b = pem.EncodeToMemory(&pem.Block{Type: "PRIVATE KEY", Bytes: b})
			os.WriteFile(path.Join(os.Getenv("APPCAST_PATH"), "dsa_key"), b, 0o600)
		}
		if test.want.Ed25519Key != nil {
			b, _ := ed25519.MarshalPrivateKey(test.want.Ed25519Key)
			b = pem.EncodeToMemory(&pem.Block{Type: "PRIVATE KEY", Bytes: b})
			os.WriteFile(path.Join(os.Getenv("APPCAST_PATH"), "ed25519_key"), b, 0o600)
		}

		got, err := config.GetSparkle(test.in)
		if err != nil {
			t.Fatal(err)
		}

		if diff := cmp.Diff(test.want, got, cmp.Exporter(func(t reflect.Type) bool { return true })); diff != "" {
			t.Error(diff)
		}
	}
}
