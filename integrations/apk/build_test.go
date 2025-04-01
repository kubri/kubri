package apk_test

import (
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"testing"

	"github.com/google/go-cmp/cmp"
	"gitlab.alpinelinux.org/alpine/go/repository"

	"github.com/kubri/kubri/integrations/apk"
	"github.com/kubri/kubri/internal/test"
	"github.com/kubri/kubri/pkg/crypto/rsa"
	source "github.com/kubri/kubri/source/file"
	target "github.com/kubri/kubri/target/file"
)

func TestBuild(t *testing.T) {
	dir := t.TempDir() + "/rpm"
	test.Golden(t, "testdata", dir, test.Ignore("*.apk"))

	want := readTestData(t)

	src, _ := source.New(source.Config{Path: "../../testdata"})
	tgt, _ := target.New(target.Config{Path: dir})

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
			KeyName: "test@example.com",
		}

		err := apk.Build(t.Context(), c)
		if err != nil {
			t.Fatal(err)
		}

		pubBytes, _ := os.ReadFile(filepath.Join(dir, "test@example.com.rsa.pub"))
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

	err := apk.Build(t.Context(), c)
	if err != nil {
		t.Fatal(err)
	}

	for name, data := range want {
		r, err := c.Target.NewReader(t.Context(), name)
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
