package apk_test

import (
	"context"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"testing"

	"github.com/abemedia/appcast/integrations/apk"
	"github.com/abemedia/appcast/internal/test"
	"github.com/abemedia/appcast/pkg/crypto/rsa"
	source "github.com/abemedia/appcast/source/file"
	target "github.com/abemedia/appcast/target/file"
	"github.com/google/go-cmp/cmp"
	"gitlab.alpinelinux.org/alpine/go/repository"
)

func TestBuild(t *testing.T) {
	want := readTestData(t)

	dir := t.TempDir() + "/rpm"
	src, _ := source.New(source.Config{Path: "../../testdata"})
	tgt, _ := target.New(target.Config{Path: dir})

	test.Golden(t, "testdata", dir)

	c := &apk.Config{
		Source: src,
		Target: tgt,
	}

	t.Run("New", func(t *testing.T) {
		testBuild(t, c, want)
	})

	t.Run("NoChange", func(t *testing.T) {
		testBuild(t, c, want)
	})

	t.Run("RSA", func(t *testing.T) {
		rsaKey, _ := rsa.NewPrivateKey()
		dir := t.TempDir()
		tgt, _ := target.New(target.Config{Path: dir})

		c := &apk.Config{
			Source:  src,
			Target:  tgt,
			RSAKey:  rsaKey,
			KeyName: "test.rsa.pub",
		}

		err := apk.Build(context.Background(), c)
		if err != nil {
			t.Fatal(err)
		}

		pubBytes, _ := os.ReadFile(filepath.Join(dir, "test.rsa.pub"))
		pub, _ := rsa.UnmarshalPublicKey(pubBytes)
		if !rsaKey.PublicKey.Equal(pub) {
			t.Fatal("should have public key")
		}

		indexFile, _ := os.Open(filepath.Join(dir, "x86_64", "APKINDEX.tar.gz"))
		apkindex, _ := repository.IndexFromArchive(indexFile)
		unsigedIndex, _ := repository.ArchiveFromIndex(apkindex)
		indexBytes, _ := io.ReadAll(unsigedIndex)
		if !rsa.Verify(pub, indexBytes, apkindex.Signature) {
			t.Fatal("should pass RSA verification")
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

func testBuild(t *testing.T, c *apk.Config, want map[string][]byte) {
	t.Helper()

	err := apk.Build(context.Background(), c)
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
