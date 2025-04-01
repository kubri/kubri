package config_test

import (
	"os"
	"testing"

	"github.com/google/go-cmp/cmp"

	"github.com/kubri/kubri/internal/test"
	"github.com/kubri/kubri/pkg/config"
)

func TestSchema(t *testing.T) {
	got := config.Schema()
	want, _ := os.ReadFile("testdata/jsonschema.json")

	if diff := cmp.Diff(want, got); diff != "" {
		t.Error(diff)
	}

	test.GoldenFile(t, "testdata/jsonschema.json", got)
}
