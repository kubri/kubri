package test

import (
	"context"
	"os"
	"reflect"
	"testing"

	"github.com/google/go-cmp/cmp"

	source "github.com/kubri/kubri/source/file"
	target "github.com/kubri/kubri/target/file"
)

func Build[C any](
	t *testing.T,
	build func(context.Context, *C) error,
	config *C,
	dir string,
) {
	t.Helper()

	src, _ := source.New(source.Config{Path: "../../testdata"})
	tgt, _ := target.New(target.Config{Path: dir})

	if config == nil {
		config = new(C)
	}
	v := reflect.ValueOf(config).Elem()
	v.FieldByName("Source").Set(reflect.ValueOf(src))
	v.FieldByName("Target").Set(reflect.ValueOf(tgt))

	t.Run("New", func(t *testing.T) {
		t.Helper()

		if err := build(t.Context(), config); err != nil {
			t.Fatal(err)
		}
		if diff := cmp.Diff(os.DirFS("testdata"), os.DirFS(dir), CompareFS()); diff != "" {
			t.Fatal(diff)
		}
	})

	t.Run("NoChange", func(t *testing.T) {
		t.Helper()

		want := ReadFS(os.DirFS(dir))
		if err := build(t.Context(), config); err != nil {
			t.Fatal(err)
		}
		got := ReadFS(os.DirFS(dir))
		if diff := cmp.Diff(want, got); diff != "" {
			t.Fatal(diff)
		}
	})
}
