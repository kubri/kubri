package slices_test

import (
	"testing"

	"github.com/abemedia/appcast/pkg/slices"
	"github.com/google/go-cmp/cmp"
)

func TestFilter(t *testing.T) {
	in := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
	want := []int{2, 4, 6, 8, 10}
	got := slices.Filter(in, func(i int) bool { return i%2 == 0 })
	if diff := cmp.Diff(want, got); diff != "" {
		t.Error(diff)
	}
}
