package cmd_test

import (
	"os"
	"testing"

	"github.com/abemedia/appcast/pkg/cmd"
)

func TestKeysCmd(t *testing.T) {
	capture(t, os.Stdout) // Reduce noise.
	err := cmd.Execute("", []string{"keys"})
	if err != nil {
		t.Error(err)
	}
}
