package arch_test

import (
	"io/fs"
	"os"
	"path"
	"testing"
	"testing/fstest"

	"github.com/google/go-cmp/cmp"

	"github.com/kubri/kubri/integrations/arch"
	"github.com/kubri/kubri/internal/test"
	"github.com/kubri/kubri/pkg/crypto/pgp"
	source "github.com/kubri/kubri/source/file"
	target "github.com/kubri/kubri/target/file"
)

func TestBuild(t *testing.T) {
	dir := t.TempDir() + "/arch"
	test.Golden(t, "testdata", dir)

	test.Build(t, arch.Build, &arch.Config{RepoName: "kubri-test"}, dir)

	t.Run("PGP", func(t *testing.T) {
		dir := t.TempDir()

		c := &arch.Config{RepoName: "kubri-test"}
		c.Source, _ = source.New(source.Config{Path: "../../testdata"})
		c.Target, _ = target.New(target.Config{Path: dir})
		c.PGPKey, _ = pgp.NewPrivateKey("test", "test@example.com")

		if err := arch.Build(t.Context(), c); err != nil {
			t.Fatal(err)
		}

		want := test.ReadFS(os.DirFS("testdata"))
		got := test.ReadFS(os.DirFS(dir))

		opt := cmp.Options{
			test.IgnoreFSMeta(),
			test.IgnoreKeys("key.asc", "*/*.sig"),
		}
		if diff := cmp.Diff(want, got, opt); diff != "" {
			t.Fatal(diff)
		}

		err := fstest.TestFS(os.DirFS(dir),
			"key.asc",
			"i686/kubri-test.db.sig",
			"i686/kubri-test-2.0.0-1-i686.pkg.tar.zst.sig",
			"x86_64/kubri-test.db.sig",
			"x86_64/kubri-test-2.0.0-1-x86_64.pkg.tar.zst.sig",
		)
		if err != nil {
			t.Fatal(err)
		}

		pub, _ := pgp.UnmarshalPublicKey(got["key.asc"].Data)
		if diff := cmp.Diff(pub, c.PGPKey, test.ComparePGPKeys()); diff != "" {
			t.Fatal(diff)
		}

		fs.WalkDir(got, ".", func(p string, d fs.DirEntry, err error) error {
			if err != nil || d.IsDir() || p == "key.asc" || path.Ext(p) == ".sig" {
				return err
			}
			if !pgp.Verify(pub, got[p].Data, got[p+".sig"].Data) {
				t.Error(p, "should pass pgp verification")
			}
			return nil
		})
	})
}
