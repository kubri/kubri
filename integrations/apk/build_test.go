package apk_test

import (
	"bytes"
	"io"
	"os"
	"testing"
	"testing/fstest"

	"gitlab.alpinelinux.org/alpine/go/repository"

	"github.com/google/go-cmp/cmp"

	"github.com/kubri/kubri/integrations/apk"
	"github.com/kubri/kubri/internal/test"
	"github.com/kubri/kubri/pkg/crypto/rsa"
	source "github.com/kubri/kubri/source/file"
	target "github.com/kubri/kubri/target/file"
)

func TestBuild(t *testing.T) {
	dir := t.TempDir()
	test.Golden(t, "testdata", dir)

	test.Build(t, apk.Build, nil, dir)

	t.Run("RSA", func(t *testing.T) {
		dir := t.TempDir()

		c := &apk.Config{KeyName: "test@example.com"}
		c.Source, _ = source.New(source.Config{Path: "../../testdata"})
		c.Target, _ = target.New(target.Config{Path: dir})
		c.RSAKey, _ = rsa.NewPrivateKey()

		if err := apk.Build(t.Context(), c); err != nil {
			t.Fatal(err)
		}

		want := test.ReadFS(os.DirFS("testdata"))
		got := test.ReadFS(os.DirFS(dir))

		opt := cmp.Options{
			test.IgnoreFSMeta(),
			test.IgnoreKeys("test@example.com.rsa.pub", "*/APKINDEX.tar.gz"),
		}
		if diff := cmp.Diff(want, got, opt); diff != "" {
			t.Fatal(diff)
		}

		err := fstest.TestFS(got, "test@example.com.rsa.pub", "x86_64/APKINDEX.tar.gz", "x86/APKINDEX.tar.gz")
		if err != nil {
			t.Fatal(err)
		}

		pub, _ := rsa.UnmarshalPublicKey(got["test@example.com.rsa.pub"].Data)
		if diff := cmp.Diff(pub, rsa.Public(c.RSAKey)); diff != "" {
			t.Fatal(diff)
		}

		for _, arch := range []string{"x86_64", "x86"} {
			wantIndex, _ := repository.IndexFromArchive(io.NopCloser(bytes.NewReader(want[arch+"/APKINDEX.tar.gz"].Data)))
			gotIndex, _ := repository.IndexFromArchive(io.NopCloser(bytes.NewReader(got[arch+"/APKINDEX.tar.gz"].Data)))

			if diff := cmp.Diff(wantIndex.Packages, gotIndex.Packages); diff != "" {
				t.Fatal(diff)
			}

			unsignedIndex, _ := repository.ArchiveFromIndex(gotIndex)
			indexBytes, _ := io.ReadAll(unsignedIndex)
			if !rsa.Verify(pub, indexBytes, gotIndex.Signature) {
				t.Fatal("should pass RSA verification")
			}
		}
	})
}
