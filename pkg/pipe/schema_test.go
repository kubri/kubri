package pipe_test

import (
	"os"
	"testing"

	"github.com/abemedia/appcast/internal/test"
	"github.com/abemedia/appcast/pkg/pipe"
	"github.com/google/go-cmp/cmp"
)

func TestSchema(t *testing.T) {
	got := pipe.Schema()
	want, _ := os.ReadFile("testdata/jsonschema.json")

	if diff := cmp.Diff(want, got); diff != "" {
		t.Error(diff)
	}

	if test.Update {
		os.WriteFile("testdata/jsonschema.json", got, 0o644)
	}
}
