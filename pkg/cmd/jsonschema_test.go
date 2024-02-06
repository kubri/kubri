package cmd_test

import (
	"bytes"
	"testing"

	"github.com/google/go-cmp/cmp"

	"github.com/kubri/kubri/pkg/cmd"
	"github.com/kubri/kubri/pkg/config"
)

func TestJsonschemaCmd(t *testing.T) {
	want := string(config.Schema())

	var stdout bytes.Buffer
	err := cmd.Execute("", cmd.WithArgs("jsonschema"), cmd.WithStdout(&stdout))
	if err != nil {
		t.Error(err)
	}

	if diff := cmp.Diff(want, stdout.String()); diff != "" {
		t.Error(diff)
	}
}
