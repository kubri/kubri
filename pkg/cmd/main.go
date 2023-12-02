package cmd

import (
	"context"
	"io"
	"log"
	"os"
	"os/signal"

	"github.com/spf13/cobra"
)

type Option func(*cobra.Command)

func WithArgs(args ...string) Option { return func(c *cobra.Command) { c.SetArgs(args) } }
func WithStdout(w io.Writer) Option  { return func(c *cobra.Command) { c.SetOut(w) } }
func WithStderr(w io.Writer) Option  { return func(c *cobra.Command) { c.SetErr(w) } }

func Execute(version string, opt ...Option) error {
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

	for _, o := range opt {
		o(cmd)
	}

	return cmd.ExecuteContext(ctx)
}

func rootCmd(version string) *cobra.Command {
	var silent bool

	cmd := &cobra.Command{
		Use:  "appcast",
		Long: "Sign and release software for common package managers and software update frameworks.",
		PersistentPreRun: func(cmd *cobra.Command, _ []string) {
			log.SetFlags(0)
			if silent {
				log.SetOutput(io.Discard)
			} else {
				log.SetOutput(cmd.OutOrStdout())
			}
		},
	}

	cmd.PersistentFlags().BoolVarP(&silent, "silent", "s", false, "only log fatal errors")

	cmd.AddCommand(buildCmd(), keysCmd(), versionCmd(version))

	return cmd
}
