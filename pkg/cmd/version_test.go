package cmd_test

import (
	"bytes"
	"fmt"
	"runtime"
	"testing"

	"github.com/google/go-cmp/cmp"

	"github.com/kubri/kubri/pkg/cmd"
)

func TestVersionCmd(t *testing.T) {
	version := "v1.0.0"
	want := fmt.Sprintf("kubri v1.0.0 %s/%s\n", runtime.GOOS, runtime.GOARCH)

	var stdout bytes.Buffer

	err := cmd.Execute(version, cmd.WithArgs("version"), cmd.WithStdout(&stdout))
	if err != nil {
		t.Error(err)
	}

	if diff := cmp.Diff(want, stdout.String()); diff != "" {
		t.Error(diff)
	}
}
