package main

import (
	"io"
	"log"
	"os"

	_ "github.com/abemedia/appcast/source/blob"
	_ "github.com/abemedia/appcast/source/github"
	_ "github.com/abemedia/appcast/source/gitlab"
	_ "github.com/abemedia/appcast/source/local"
	"github.com/spf13/cobra"
)

func main() {
	var verbose bool

	cmd := &cobra.Command{
		Use:  "appcast",
		Long: "Generate appcast XML files for Sparkle from your repo.",
		PersistentPreRun: func(*cobra.Command, []string) {
			if verbose {
				log.SetFlags(0)
			} else {
				log.SetOutput(io.Discard)
			}
		},
	}

	cmd.PersistentFlags().BoolVarP(&verbose, "verbose", "v", false, "enable verbose logging")

	cmd.AddCommand(
		feedCmd(),
		signCmd(),
		keysCmd(),
		versionCmd(),
	)

	if err := cmd.Execute(); err != nil {
		os.Exit(1)
	}
}
