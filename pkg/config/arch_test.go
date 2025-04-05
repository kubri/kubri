package config_test

import (
	"fmt"
	"testing"

	"github.com/kubri/kubri/integrations/arch"
	"github.com/kubri/kubri/pkg/config"
	"github.com/kubri/kubri/pkg/crypto"
	"github.com/kubri/kubri/pkg/crypto/pgp"
	"github.com/kubri/kubri/pkg/secret"
	source "github.com/kubri/kubri/source/file"
	target "github.com/kubri/kubri/target/file"
)

func TestArch(t *testing.T) {
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
				arch:
					disabled: true
					repo-name: kubri-test
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
				arch:
					repo-name: kubri-test
			`,
			want: &config.Config{
				Arch: &arch.Config{
					RepoName: "kubri-test",
					Source:   src,
					Target:   tgt.Sub("arch"),
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
				arch:
					folder: test
					repo-name: kubri-test
			`,
			hook: func() { secret.Put("pgp_key", keyBytes) },
			want: &config.Config{
				Arch: &arch.Config{
					RepoName:   "kubri-test",
					Source:     src,
					Target:     tgt.Sub("test"),
					Version:    "latest",
					Prerelease: true,
					PGPKey:     key,
				},
			},
		},
		{
			desc: "missing repo name",
			in: `
				source:
					type: file
					path: ` + dir + `
				target:
					type: file
					path: ` + dir + `
				arch: {}
			`,
			err: &config.Error{Errors: []string{"arch.repo-name is a required field"}},
		},
		{
			desc: "invalid repo name",
			in: `
				source:
					type: file
					path: ` + dir + `
				target:
					type: file
					path: ` + dir + `
				arch:
					repo-name: 'Foo123 Bar'
			`,
			err: &config.Error{Errors: []string{"arch.repo-name must only contain letters, numbers, dashes and underscores"}},
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
				arch:
					repo-name: kubri-test
			`,
			hook: func() { secret.Put("pgp_key", []byte("nope")) },
			err:  fmt.Errorf("%w: no armored data found", crypto.ErrInvalidKey),
		},
		{
			desc: "invalid folder",
			in: `
				version: latest
				prerelease: true
				source:
					type: file
					path: ` + dir + `
				target:
					type: file
					path: ` + dir + `
				arch:
					folder: '*'
					repo-name: kubri-test
			`,
			err: &config.Error{Errors: []string{"arch.folder must be a valid folder name"}},
		},
		{
			desc: "absolute folder",
			in: `
				version: latest
				prerelease: true
				source:
					type: file
					path: ` + dir + `
				target:
					type: file
					path: ` + dir + `
				arch:
					folder: '/foo/bar'
					repo-name: kubri-test
			`,
			err: &config.Error{Errors: []string{"arch.folder must be a valid folder name"}},
		},
	})
}
