package config_test

import (
	"fmt"
	"testing"

	"github.com/kubri/kubri/integrations/apt"
	"github.com/kubri/kubri/pkg/config"
	"github.com/kubri/kubri/pkg/crypto"
	"github.com/kubri/kubri/pkg/crypto/pgp"
	"github.com/kubri/kubri/pkg/secret"
	source "github.com/kubri/kubri/source/file"
	target "github.com/kubri/kubri/target/file"
)

func TestApt(t *testing.T) {
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
				apt:
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
				apt: {}
			`,
			want: &config.Config{
				Apt: &apt.Config{
					Source: src,
					Target: tgt.Sub("apt"),
				},
			},
		},
		{
			desc: "no compression",
			in: `
				source:
					type: file
					path: ` + dir + `
				target:
					type: file
					path: ` + dir + `
				apt:
					compress: [none]
			`,
			want: &config.Config{
				Apt: &apt.Config{
					Source:   src,
					Target:   tgt.Sub("apt"),
					Compress: apt.NoCompression,
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
				apt:
					folder: test
					compress:
						- gzip
						- bzip2
						- xz
						- lzma
						- lz4
						- zstd
			`,
			hook: func() { secret.Put("pgp_key", keyBytes) },
			want: &config.Config{
				Apt: &apt.Config{
					Source:     src,
					Target:     tgt.Sub("test"),
					Version:    "latest",
					Prerelease: true,
					PGPKey:     key,
					Compress:   apt.GZIP | apt.BZIP2 | apt.XZ | apt.LZMA | apt.LZ4 | apt.ZSTD,
				},
			},
		},
		{
			desc: "validation",
			in: `
				source:
					type: file
					path: ` + dir + `
				target:
					type: file
					path: ` + dir + `
				apt:
					folder: '*'
					compress: [invalid]
			`,
			err: &config.Error{
				Errors: []string{
					"apt.folder must be a valid folder name",
					"apt.compress[0] must be one of [none gzip bzip2 xz lzma lz4 zstd]",
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
				apt: {}
			`,
			hook: func() { secret.Put("pgp_key", []byte("nope")) },
			err:  fmt.Errorf("%w: no armored data found", crypto.ErrInvalidKey),
		},
	})
}
