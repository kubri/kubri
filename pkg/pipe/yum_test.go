package pipe_test

import (
	"fmt"
	"testing"

	"github.com/abemedia/appcast/integrations/yum"
	"github.com/abemedia/appcast/pkg/crypto"
	"github.com/abemedia/appcast/pkg/crypto/pgp"
	"github.com/abemedia/appcast/pkg/pipe"
	"github.com/abemedia/appcast/pkg/secret"
	source "github.com/abemedia/appcast/source/file"
	target "github.com/abemedia/appcast/target/file"
)

func TestYum(t *testing.T) {
	dir := t.TempDir()
	src, _ := source.New(source.Config{Path: dir})
	tgt, _ := target.New(target.Config{Path: dir})
	key, _ := pgp.NewPrivateKey("test", "test@example.com")
	keyBytes, _ := pgp.MarshalPrivateKey(key)

	runTest(t, []testCase{
		{
			desc: "disabled",
			in: `
				source:
					type: file
					path: ` + dir + `
				target:
					type: file
					path: ` + dir + `
				yum:
					disabled: true
			`,
			want: &pipe.Pipe{},
		},
		{
			desc: "defaults",
			in: `
				source:
					type: file
					path: ` + dir + `
				target:
					type: file
					path: ` + dir + `
				yum: {}
			`,
			want: &pipe.Pipe{
				Yum: &yum.Config{
					Source: src,
					Target: tgt.Sub("yum"),
				},
			},
		},
		{
			desc: "full",
			in: `
				version: latest
				prerelease: true
				source:
					type: file
					path: ` + dir + `
				target:
					type: file
					path: ` + dir + `
				yum:
					folder: test
			`,
			hook: func() { secret.Put("pgp_key", keyBytes) },
			want: &pipe.Pipe{
				Yum: &yum.Config{
					Source:     src,
					Target:     tgt.Sub("test"),
					Version:    "latest",
					Prerelease: true,
					PGPKey:     key,
				},
			},
		},
		{
			desc: "invalid pgp key",
			in: `
				source:
					type: file
					path: ` + dir + `
				target:
					type: file
					path: ` + dir + `
				yum: {}
			`,
			hook: func() { secret.Put("pgp_key", []byte("nope")) },
			err:  fmt.Errorf("%w: no armored data found", crypto.ErrInvalidKey),
		},
	})
}