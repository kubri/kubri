package config_test

import (
	"errors"
	"testing"

	"github.com/abemedia/appcast/integrations/apk"
	"github.com/abemedia/appcast/pkg/config"
	"github.com/abemedia/appcast/pkg/crypto"
	"github.com/abemedia/appcast/pkg/crypto/rsa"
	"github.com/abemedia/appcast/pkg/secret"
	source "github.com/abemedia/appcast/source/file"
	target "github.com/abemedia/appcast/target/file"
)

func TestApk(t *testing.T) {
	dir := t.TempDir()
	src, _ := source.New(source.Config{Path: dir})
	tgt, _ := target.New(target.Config{Path: dir})
	key, _ := rsa.NewPrivateKey()
	keyBytes, _ := rsa.MarshalPrivateKey(key)

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
				apk:
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
				apk: {}
			`,
			want: &config.Config{
				Apk: &apk.Config{
					Source: src,
					Target: tgt.Sub("apk"),
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
				apk:
					folder: test
					key-name: test@example.com.rsa.pub
			`,
			hook: func() { secret.Put("rsa_key", keyBytes) },
			want: &config.Config{
				Apk: &apk.Config{
					Source:     src,
					Target:     tgt.Sub("test"),
					Version:    "latest",
					Prerelease: true,
					RSAKey:     key,
					KeyName:    "test@example.com.rsa.pub",
				},
			},
		},
		{
			desc: "missing key name",
			in: `
				source:
					type: file
					path: ` + dir + `
				target:
					type: file
					path: ` + dir + `
				apk: {}
			`,
			hook: func() { secret.Put("rsa_key", keyBytes) },
			err:  errors.New("missing key name"),
		},
		{
			desc: "invalid rsa key",
			in: `
				source:
					type: file
					path: ` + dir + `
				target:
					type: file
					path: ` + dir + `
				apk: {}
			`,
			hook: func() { secret.Put("rsa_key", []byte("nope")) },
			err:  crypto.ErrInvalidKey,
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
				apk:
					folder: '*'
			`,
			err: &config.Error{Errors: []string{"apk.folder must be a valid folder name"}},
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
				apk:
					folder: '/foo/bar'
			`,
			err: &config.Error{Errors: []string{"apk.folder must be a valid folder name"}},
		},
	})
}
