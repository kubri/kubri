package cmd

import (
	"github.com/abemedia/appcast/pkg/pipe"
	"github.com/spf13/cobra"
)

func buildCmd() *cobra.Command {
	var configPath string

	cmd := &cobra.Command{
		Use:     "build",
		Short:   "Build appcast feed",
		Aliases: []string{"b"},
		RunE: func(cmd *cobra.Command, args []string) error {
			p, err := pipe.Load(configPath)
			if err != nil {
				return err
			}
			return p.Run(cmd.Context())
		},
	}

	cmd.Flags().StringVarP(&configPath, "config", "c", "", "load configuration from a file")

	return cmd
}
