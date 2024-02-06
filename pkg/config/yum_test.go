package config_test

import (
	"fmt"
	"testing"

	"github.com/kubri/kubri/integrations/yum"
	"github.com/kubri/kubri/pkg/config"
	"github.com/kubri/kubri/pkg/crypto"
	"github.com/kubri/kubri/pkg/crypto/pgp"
	"github.com/kubri/kubri/pkg/secret"
	source "github.com/kubri/kubri/source/file"
	target "github.com/kubri/kubri/target/file"
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
			want: &config.Config{},
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
			want: &config.Config{
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
			want: &config.Config{
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
		{
			desc: "invalid folder",
			in: `
				source:
					type: file
					path: ` + dir + `
				target:
					type: file
					path: ` + dir + `
				yum:
					folder: '*'
			`,
			err: &config.Error{Errors: []string{"yum.folder must be a valid folder name"}},
		},
	})
}
