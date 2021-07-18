package main

import (
	"runtime"

	"github.com/spf13/cobra"
)

var version = "master"

func versionCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "version",
		Short:   "Print current version",
		Aliases: []string{"s"},
		Args:    cobra.NoArgs,
		Run: func(cmd *cobra.Command, _ []string) {
			cmd.PrintErrf("appcast %s %s/%s\n", version, runtime.GOOS, runtime.GOARCH)
		},
	}

	return cmd
}
