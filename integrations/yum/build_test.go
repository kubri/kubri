package yum_test

import (
	"bytes"
	"compress/gzip"
	"context"
	"io"
	"io/fs"
	"os"
	"path"
	"strconv"
	"testing"
	"time"

	"github.com/abemedia/appcast/integrations/yum"
	source "github.com/abemedia/appcast/source/file"
	target "github.com/abemedia/appcast/target/file"
	"github.com/google/go-cmp/cmp"
)

func TestBuild(t *testing.T) {
	want := readTestData(t)
	dir := t.TempDir() + "/rpm"

	src, err := source.New(source.Config{Path: "../../testdata"})
	if err != nil {
		t.Fatal(err)
	}

	tgt, err := target.New(target.Config{Path: dir})
	if err != nil {
		t.Fatal(err)
	}

	c := &yum.Config{
		Source: src,
		Target: tgt,
	}

	testBuild(t, c, os.DirFS(dir), want)

	// Should be no-op as nothing changed so timestamps should still be valid.
	time.Sleep(time.Second)
	testBuild(t, c, os.DirFS(dir), want)
}

func readTestData(t *testing.T) map[string]string {
	t.Helper()

	ts := strconv.Itoa(int(time.Now().Unix()))
	want := make(map[string]string)

	err := fs.WalkDir(os.DirFS("testdata"), ".", func(path string, d fs.DirEntry, err error) error {
		if err != nil || d.IsDir() {
			return err
		}
		b, err := fs.ReadFile(os.DirFS("testdata"), d.Name())
		if err != nil {
			return err
		}
		b = bytes.TrimSpace(bytes.ReplaceAll(b, []byte("__TIME__"), []byte(ts)))
		path = "*-" + path + ".gz"
		want["repodata/"+path] = string(b)
		return nil
	})
	if err != nil {
		t.Fatal(err)
	}

	return want
}

func testBuild(t *testing.T, c *yum.Config, fsys fs.FS, want map[string]string) {
	t.Helper()

	err := yum.Build(context.Background(), c)
	if err != nil {
		t.Fatal(err)
	}

	for name, data := range want {
		matches, err := fs.Glob(fsys, name)
		if err != nil {
			t.Fatal(name, err)
		}
		if len(matches) != 1 {
			t.Fatalf("Expected %s: %v", name, matches)
		}
		r, err := c.Target.NewReader(context.Background(), matches[0])
		if err != nil {
			t.Fatal(name, err)
		}
		defer r.Close()

		if path.Ext(name) == ".gz" {
			r, err = gzip.NewReader(r)
			if err != nil {
				t.Fatal(name, err)
			}
		}

		got, err := io.ReadAll(r)
		if err != nil {
			t.Fatal(name, err)
		}

		if diff := cmp.Diff(data, string(got)); diff != "" {
			t.Error(name, diff)
		}
	}
}
