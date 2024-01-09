package cmd_test

import (
	"bytes"
	"testing"

	"github.com/abemedia/appcast/pkg/cmd"
	"github.com/abemedia/appcast/pkg/pipe"
	"github.com/google/go-cmp/cmp"
)

func TestJsonschemaCmd(t *testing.T) {
	want := string(pipe.Schema())

	var stdout bytes.Buffer
	err := cmd.Execute("", cmd.WithArgs("jsonschema"), cmd.WithStdout(&stdout))
	if err != nil {
		t.Error(err)
	}

	if diff := cmp.Diff(want, stdout.String()); diff != "" {
		t.Error(diff)
	}
}
