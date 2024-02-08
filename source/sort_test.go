package source_test

import (
	"sort"
	"testing"

	"github.com/google/go-cmp/cmp"

	"github.com/kubri/kubri/source"
)

func TestByVersion(t *testing.T) {
	in := []*source.Release{
		{Version: "0.9.1"},
		{Version: "1.0.0"},
		{Version: "1.51.0"},
		{Version: "1.5.9"},
	}
	want := []*source.Release{
		{Version: "1.51.0"},
		{Version: "1.5.9"},
		{Version: "1.0.0"},
		{Version: "0.9.1"},
	}

	sort.Sort(source.ByVersion(in))

	if diff := cmp.Diff(want, in); diff != "" {
		t.Error(diff)
	}
}
