package main

import (
	"log"
	"runtime"

	"github.com/spf13/cobra"
)

var version string

func versionCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "version",
		Short:   "Print current version",
		Aliases: []string{"s"},
		Args:    cobra.NoArgs,
		Run: func(*cobra.Command, []string) {
			log.Printf("appcast %s %s/%s", version, runtime.GOOS, runtime.GOARCH)
		},
	}

	return cmd
}
