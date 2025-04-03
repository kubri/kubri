package yum_test

import (
	"os"
	"testing"
	"testing/fstest"
	"time"

	"github.com/google/go-cmp/cmp"

	"github.com/kubri/kubri/integrations/yum"
	"github.com/kubri/kubri/internal/test"
	"github.com/kubri/kubri/pkg/crypto/pgp"
	source "github.com/kubri/kubri/source/file"
	target "github.com/kubri/kubri/target/file"
)

func TestBuild(t *testing.T) {
	dir := t.TempDir()
	test.Golden(t, "testdata", dir)

	yum.SetTime(time.Date(2023, 11, 19, 23, 37, 12, 0, time.UTC))

	test.Build(t, yum.Build, nil, dir)

	t.Run("PGP", func(t *testing.T) {
		dir := t.TempDir()

		c := &yum.Config{}
		c.Source, _ = source.New(source.Config{Path: "../../testdata"})
		c.Target, _ = target.New(target.Config{Path: dir})
		c.PGPKey, _ = pgp.NewPrivateKey("test", "test@example.com")

		if err := yum.Build(t.Context(), c); err != nil {
			t.Fatal(err)
		}

		want := test.ReadFS(os.DirFS("testdata"))
		got := test.ReadFS(os.DirFS(dir))

		opt := cmp.Options{
			test.IgnoreFSMeta(),
			test.IgnoreKeys("repodata/repomd.xml.asc", "repodata/repomd.xml.key"),
		}
		if diff := cmp.Diff(want, got, opt); diff != "" {
			t.Fatal(diff)
		}

		err := fstest.TestFS(got, "repodata/repomd.xml.asc", "repodata/repomd.xml.key")
		if err != nil {
			t.Fatal(err)
		}

		pub, _ := pgp.UnmarshalPublicKey(got["repodata/repomd.xml.key"].Data)
		if diff := cmp.Diff(pub, c.PGPKey, test.ComparePGPKeys()); diff != "" {
			t.Fatal(diff)
		}
		if !pgp.Verify(pub, got["repodata/repomd.xml"].Data, got["repodata/repomd.xml.asc"].Data) {
			t.Error("should pass pgp verification")
		}
	})
}
