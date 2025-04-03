package apt_test

import (
	"os"
	"testing"
	"testing/fstest"
	"time"

	"github.com/google/go-cmp/cmp"

	"github.com/kubri/kubri/integrations/apt"
	"github.com/kubri/kubri/internal/test"
	"github.com/kubri/kubri/pkg/crypto/pgp"
	source "github.com/kubri/kubri/source/file"
	target "github.com/kubri/kubri/target/file"
)

func TestBuild(t *testing.T) {
	dir := t.TempDir()
	test.Golden(t, "testdata", dir)

	apt.SetTime(time.Date(2023, 11, 19, 23, 37, 12, 0, time.UTC))

	test.Build(t, apt.Build, nil, dir)

	t.Run("PGP", func(t *testing.T) {
		dir := t.TempDir()

		c := &apt.Config{}
		c.Source, _ = source.New(source.Config{Path: "../../testdata"})
		c.Target, _ = target.New(target.Config{Path: dir})
		c.PGPKey, _ = pgp.NewPrivateKey("test", "test@example.com")

		if err := apt.Build(t.Context(), c); err != nil {
			t.Fatal(err)
		}

		want := test.ReadFS(os.DirFS("testdata"))
		got := test.ReadFS(os.DirFS(dir))

		opt := cmp.Options{
			test.IgnoreFSMeta(),
			test.IgnoreKeys("key.asc", "dists/stable/Release.gpg", "dists/stable/InRelease"),
		}
		if diff := cmp.Diff(want, got, opt); diff != "" {
			t.Fatal(diff)
		}

		err := fstest.TestFS(os.DirFS(dir), "key.asc", "dists/stable/Release.gpg", "dists/stable/InRelease")
		if err != nil {
			t.Fatal(err)
		}

		pub, _ := pgp.UnmarshalPublicKey(got["key.asc"].Data)
		if diff := cmp.Diff(pub, c.PGPKey, test.ComparePGPKeys()); diff != "" {
			t.Fatal(diff)
		}
		if !pgp.Verify(pub, got["dists/stable/Release"].Data, got["dists/stable/Release.gpg"].Data) {
			t.Error("should pass pgp verification")
		}

		data, sig, err := pgp.Split(got["dists/stable/InRelease"].Data)
		if err != nil {
			t.Fatal(err)
		}
		if diff := cmp.Diff(want["dists/stable/Release"].Data, data); diff != "" {
			t.Error(diff)
		}
		if !pgp.Verify(pub, data, sig) {
			t.Error("should pass pgp verification")
		}
	})

	t.Run("CustomCompress", func(t *testing.T) {
		dir := t.TempDir()

		c := &apt.Config{Compress: apt.BZIP2 | apt.ZSTD}
		c.Source, _ = source.New(source.Config{Path: "../../testdata"})
		c.Target, _ = target.New(target.Config{Path: dir})

		if err := apt.Build(t.Context(), c); err != nil {
			t.Fatal(err)
		}

		err := fstest.TestFS(os.DirFS(dir),
			"dists/stable/main/binary-amd64/Packages",
			"dists/stable/main/binary-amd64/Packages.bz2",
			"dists/stable/main/binary-amd64/Packages.zst",
			"dists/stable/main/binary-i386/Packages",
			"dists/stable/main/binary-i386/Packages.bz2",
			"dists/stable/main/binary-i386/Packages.zst",
		)
		if err != nil {
			t.Error(err)
		}

		err = fstest.TestFS(os.DirFS(dir),
			"dists/stable/main/binary-amd64/Packages.gz",
			"dists/stable/main/binary-i386/Packages.gz",
		)
		if err == nil {
			t.Error("should not have gzip files")
		}
	})
}
