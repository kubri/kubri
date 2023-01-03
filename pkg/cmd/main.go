package cmd

import (
	"context"
	"io"
	"log"
	"os"
	"os/signal"

	"github.com/spf13/cobra"
)

func Execute(version string, args []string) error {
	ctx, cancel := context.WithCancel(context.Background())
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	go func() {
		<-c // Close gracefully if CTRL+C is pressed once.
		cancel()
		<-c // Exit if CTRL+C is pressed twice.
		os.Exit(1)
	}()

	cmd := rootCmd(version)
	if args != nil {
		cmd.SetArgs(args)
	}

	return cmd.ExecuteContext(ctx)
}

func rootCmd(version string) *cobra.Command {
	var silent bool

	cmd := &cobra.Command{
		Use:  "appcast",
		Long: "Generate appcast XML files for Sparkle from your repo.",
		PersistentPreRun: func(*cobra.Command, []string) {
			log.SetFlags(0)
			if silent {
				log.SetOutput(io.Discard)
			}
		},
	}

	cmd.PersistentFlags().BoolVarP(&silent, "silent", "s", false, "only log fatal errors")

	cmd.AddCommand(buildCmd(), keysCmd(), versionCmd(version))

	return cmd
}
