package cmd

import (
	"fmt"
	"runtime"

	"github.com/spf13/cobra"
)

func versionCmd(version string) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "version",
		Short: "Print current version",
		Args:  cobra.NoArgs,
		Run: func(cmd *cobra.Command, _ []string) {
			fmt.Fprintf(cmd.OutOrStdout(), "appcast %s %s/%s\n", version, runtime.GOOS, runtime.GOARCH)
		},
	}

	return cmd
}
