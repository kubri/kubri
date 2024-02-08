package apt_test

import (
	"context"
	"io"
	"io/fs"
	"maps"
	"os"
	"path"
	"path/filepath"
	"strings"
	"testing"
	"testing/fstest"
	"time"

	"github.com/google/go-cmp/cmp"

	"github.com/kubri/kubri/integrations/apt"
	"github.com/kubri/kubri/internal/test"
	"github.com/kubri/kubri/pkg/crypto/pgp"
	source "github.com/kubri/kubri/source/file"
	"github.com/kubri/kubri/target"
	ftarget "github.com/kubri/kubri/target/file"
)

func TestBuild(t *testing.T) {
	want := readTestData(t, ".gz", ".xz")
	now := time.Date(2023, 11, 19, 23, 37, 12, 0, time.UTC)

	dir := t.TempDir() + "/apt"
	src, _ := source.New(source.Config{Path: "../../testdata"})
	tgt, _ := ftarget.New(ftarget.Config{Path: dir})

	test.Golden(t, "testdata", dir, test.Ignore("*.deb", "*.gz", "*.xz"))

	t.Run("New", func(t *testing.T) {
		c := &apt.Config{Source: src, Target: tgt}
		testBuild(t, c, want, now)
	})

	// Should be no-op as nothing changed so timestamp should still be valid.
	t.Run("NoChange", func(t *testing.T) {
		c := &apt.Config{Source: src, Target: tgt}
		testBuild(t, c, want, now.Add(time.Hour))
	})

	t.Run("PGP", func(t *testing.T) {
		dir := t.TempDir()
		pgpKey, _ := pgp.NewPrivateKey("test", "test@example.com")
		tgt, _ := ftarget.New(ftarget.Config{Path: dir})

		c := &apt.Config{
			Source: src,
			Target: tgt,
			PGPKey: pgpKey,
		}

		// Remove InRelease and test that separately below.
		wantPGP := maps.Clone(want)
		delete(wantPGP, "dists/stable/InRelease")

		testBuild(t, c, wantPGP, now)

		data, _ := os.ReadFile(filepath.Join(dir, "dists", "stable", "Release"))
		sig, _ := os.ReadFile(filepath.Join(dir, "dists", "stable", "Release.gpg"))
		if !pgp.Verify(pgp.Public(pgpKey), data, sig) {
			t.Error("should pass pgp verification")
		}

		in, _ := os.ReadFile(filepath.Join(dir, "dists", "stable", "InRelease"))
		data, sig, err := pgp.Split(in)
		if err != nil {
			t.Fatal(err)
		}
		if !pgp.Verify(pgp.Public(c.PGPKey), data, sig) {
			t.Error("failed to verify InRelease signature")
		}
		if diff := cmp.Diff(want["dists/stable/Release"], string(data)); diff != "" {
			t.Error("dists/stable/InRelease", diff)
		}
	})

	t.Run("CustomCompress", func(t *testing.T) {
		dir := t.TempDir()
		tgt, _ := ftarget.New(ftarget.Config{Path: dir})

		c := &apt.Config{
			Source:   src,
			Target:   tgt,
			Compress: apt.BZIP2 | apt.ZSTD,
		}

		err := apt.Build(context.Background(), c)
		if err != nil {
			t.Fatal(err)
		}

		err = fstest.TestFS(os.DirFS(dir),
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

func readTestData(t *testing.T, compress ...string) map[string]string {
	t.Helper()

	want := make(map[string]string)

	err := fs.WalkDir(os.DirFS("testdata"), ".", func(path string, d fs.DirEntry, err error) error {
		if err != nil || d.IsDir() {
			return err
		}
		b, err := fs.ReadFile(os.DirFS("testdata"), path)
		if err != nil {
			return err
		}
		want[path] = string(b)
		if d.Name() == "Packages" {
			for _, ext := range compress {
				want[path+ext] = string(b)
			}
		}
		return nil
	})
	if err != nil {
		t.Fatal(err)
	}

	return want
}

func testBuild(t *testing.T, c *apt.Config, want map[string]string, now time.Time) { //nolint:thelper
	apt.SetTime(now)

	err := apt.Build(context.Background(), c)
	if err != nil {
		t.Fatal(err)
	}

	for name, data := range want {
		got := readFile(t, c.Target, name)

		ext := path.Ext(name)
		base := strings.TrimSuffix(path.Base(name), ext)
		if base == "Packages" {
			data = want[strings.TrimSuffix(name, ext)]
		}

		if diff := cmp.Diff(data, string(got)); diff != "" {
			t.Error(name, diff)
		}
	}
}

func readFile(t *testing.T, tgt target.Target, name string) []byte {
	t.Helper()

	r, err := tgt.NewReader(context.Background(), name)
	if err != nil {
		t.Fatal(name, err)
	}
	defer r.Close()

	r, err = apt.Decompress(path.Ext(name))(r)
	if err != nil {
		t.Fatal(name, err)
	}
	defer r.Close()

	b, err := io.ReadAll(r)
	if err != nil {
		t.Fatal(name, err)
	}

	return b
}
