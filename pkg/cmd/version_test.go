package cmd_test

import (
	"fmt"
	"os"
	"runtime"
	"testing"

	"github.com/abemedia/appcast/pkg/cmd"
	"github.com/google/go-cmp/cmp"
)

func TestVersionCmd(t *testing.T) {
	version := "v1.0.0"
	want := fmt.Sprintf("appcast v1.0.0 %s/%s\n", runtime.GOOS, runtime.GOARCH)

	stdout := capture(t, os.Stdout)

	err := cmd.Execute(version, []string{"version"})
	if err != nil {
		t.Error(err)
	}

	if diff := cmp.Diff(want, stdout.String()); diff != "" {
		t.Error(diff)
	}
}
