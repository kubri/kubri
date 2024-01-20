package yum_test

import (
	"context"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/abemedia/appcast/integrations/yum"
	"github.com/abemedia/appcast/internal/test"
	"github.com/abemedia/appcast/pkg/crypto/pgp"
	source "github.com/abemedia/appcast/source/file"
	target "github.com/abemedia/appcast/target/file"
	"github.com/google/go-cmp/cmp"
)

func TestBuild(t *testing.T) {
	want := readTestData(t)
	now := time.Date(2023, 11, 19, 23, 37, 12, 0, time.UTC)

	dir := t.TempDir() + "/rpm"
	src, _ := source.New(source.Config{Path: "../../testdata"})
	tgt, _ := target.New(target.Config{Path: dir})
	key, _ := pgp.NewPrivateKey("test", "test@example.com")

	test.Golden(t, "testdata", dir)

	c := &yum.Config{
		Source: src,
		Target: tgt,
		PGPKey: key,
	}

	t.Run("New", func(t *testing.T) {
		testBuild(t, c, want, now)
	})

	t.Run("NoChange", func(t *testing.T) {
		testBuild(t, c, want, now.Add(time.Hour))
	})

	t.Run("PGP", func(t *testing.T) {
		dir := t.TempDir()
		pgpKey, _ := pgp.NewPrivateKey("test", "test@example.com")
		tgt, _ := target.New(target.Config{Path: dir})

		c := &yum.Config{
			Source: src,
			Target: tgt,
			PGPKey: pgpKey,
		}

		testBuild(t, c, want, now)

		data, _ := os.ReadFile(filepath.Join(dir, "repodata", "repomd.xml"))
		key, _ := os.ReadFile(filepath.Join(dir, "repodata", "repomd.xml.key"))
		sig, _ := os.ReadFile(filepath.Join(dir, "repodata", "repomd.xml.asc"))
		pub, _ := pgp.UnmarshalPublicKey(key)

		if !pgp.Verify(pub, data, sig) {
			t.Error("should pass pgp verification")
		}
	})
}

func readTestData(t *testing.T) map[string][]byte {
	t.Helper()

	want := make(map[string][]byte)

	err := fs.WalkDir(os.DirFS("testdata"), ".", func(path string, d fs.DirEntry, err error) error {
		if err != nil || d.IsDir() {
			return err
		}
		b, err := fs.ReadFile(os.DirFS("testdata"), path)
		if err != nil {
			return err
		}
		want[path] = b
		return nil
	})
	if err != nil {
		t.Fatal(err)
	}

	return want
}

func testBuild(t *testing.T, c *yum.Config, want map[string][]byte, now time.Time) {
	t.Helper()

	yum.SetTime(now)

	err := yum.Build(context.Background(), c)
	if err != nil {
		t.Fatal(err)
	}

	for name, data := range want {
		r, err := c.Target.NewReader(context.Background(), name)
		if err != nil {
			t.Fatal(name, err)
		}
		defer r.Close()

		got, err := io.ReadAll(r)
		if err != nil {
			t.Fatal(name, err)
		}

		if diff := cmp.Diff(data, got); diff != "" {
			t.Error(name, diff)
		}
	}
}
