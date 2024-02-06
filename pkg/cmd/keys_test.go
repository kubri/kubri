package cmd_test

import (
	"io"
	"testing"

	"github.com/kubri/kubri/pkg/cmd"
)

func TestKeysCmd(t *testing.T) {
	err := cmd.Execute("", cmd.WithArgs("keys"), cmd.WithStdout(io.Discard))
	if err != nil {
		t.Error(err)
	}
}
