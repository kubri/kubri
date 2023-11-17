package pipe //nolint:testpackage

import (
	"errors"
	"testing"

	"github.com/abemedia/appcast/integrations/apt"
	"github.com/abemedia/appcast/internal/test"
	"github.com/abemedia/appcast/pkg/crypto/pgp"
	source "github.com/abemedia/appcast/source/file"
	target "github.com/abemedia/appcast/target/file"
	"github.com/google/go-cmp/cmp"
)

func TestApt(t *testing.T) {
	t.Setenv("APPCAST_PATH", t.TempDir())
	src, _ := source.New(source.Config{Path: t.TempDir()})
	tgt, _ := target.New(target.Config{Path: t.TempDir()})
	key, _ := pgp.NewPrivateKey("test", "test@example.com")
	keyBytes, _ := pgp.MarshalPrivateKey(key)

	tests := []struct {
		in   *config
		want *apt.Config
		err  error
		hook func()
	}{
		{
			in: &config{
				source: src,
				target: tgt,
				Apt:    &aptConfig{},
			},
			want: &apt.Config{
				Source: src,
				Target: tgt.Sub("apt"),
			},
		},
		{
			in: &config{
				source: src,
				target: tgt,
				Apt: &aptConfig{
					Compress: []string{"none"},
				},
			},
			want: &apt.Config{
				Source:   src,
				Target:   tgt.Sub("apt"),
				Compress: apt.NoCompression,
			},
		},
		{
			in: &config{
				source: src,
				target: tgt,
				Apt: &aptConfig{
					Folder:   "deb",
					Compress: []string{"gzip", "bzip2", "xz", "lzma", "lz4", "zstd"},
				},
			},
			want: &apt.Config{
				Source:   src,
				Target:   tgt.Sub("deb"),
				PGPKey:   key,
				Compress: apt.GZIP | apt.BZIP2 | apt.XZ | apt.LZMA | apt.LZ4 | apt.ZSTD,
			},
			hook: func() { t.Setenv("APPCAST_PGP_KEY", string(keyBytes)) },
		},
		{
			in: &config{
				source: src,
				target: tgt,
				Apt: &aptConfig{
					Compress: []string{"foo"},
				},
			},
			err: errors.New("unknown compression algorithm: foo"),
		},
		{
			in: &config{
				source: src,
				target: tgt,
				Apt:    &aptConfig{},
			},
			err:  errors.New("invalid key: no armored data found"),
			hook: func() { t.Setenv("APPCAST_PGP_KEY", "nope") },
		},
	}

	opts := cmp.Options{
		test.ExportAll(),
		test.ComparePGPKeys(),
		test.CompareErrorMessages(),
	}

	for i, test := range tests {
		if test.hook != nil {
			test.hook()
		}

		s, err := getApt(test.in)

		if diff := cmp.Diff(test.err, err, opts); diff != "" {
			t.Error(i, diff)
		} else if diff := cmp.Diff(test.want, s, opts); diff != "" {
			t.Error(i, diff)
		}
	}
}
