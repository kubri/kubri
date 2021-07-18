package main

import (
	"log"
	"os"

	_ "github.com/abemedia/appcast/source/file"
	_ "github.com/abemedia/appcast/source/github"
	_ "github.com/abemedia/appcast/source/gitlab"
	_ "github.com/abemedia/appcast/source/local"
	"github.com/spf13/cobra"
)

func main() {
	log.SetFlags(0)

	cmd := &cobra.Command{
		Use:  "appcast",
		Long: "Generate appcast XML files for Sparkle from your repo.",
	}

	cmd.AddCommand(
		generateCmd(),
		signCmd(),
		keysCmd(),
	)

	if err := cmd.Execute(); err != nil {
		os.Exit(1)
	}
}
