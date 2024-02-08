package slices_test

import (
	"testing"

	"github.com/google/go-cmp/cmp"

	"github.com/kubri/kubri/pkg/slices"
)

func TestFilter(t *testing.T) {
	in := []string{"foo", "bar", "baz"}
	want := []string{"foo", "baz"}
	got := slices.Filter(in, func(e string) bool { return e != "bar" })
	if diff := cmp.Diff(want, got); diff != "" {
		t.Error(diff)
	}
}
