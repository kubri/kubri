package cmd

import (
	"github.com/abemedia/appcast/integrations/sparkle"
	"github.com/abemedia/appcast/pkg/config"
	"github.com/spf13/cobra"

	// Import source & target providers.
	_ "github.com/abemedia/appcast/source/blob"
	_ "github.com/abemedia/appcast/source/github"
	_ "github.com/abemedia/appcast/source/gitlab"
	_ "github.com/abemedia/appcast/source/local"
	_ "github.com/abemedia/appcast/target/blob"
)

func buildCmd() *cobra.Command {
	var configPath string

	cmd := &cobra.Command{
		Use:     "build",
		Short:   "Build appcast feed",
		Aliases: []string{"b"},
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := config.Load(configPath)
			if err != nil {
				return err
			}

			sparkleConfig, err := config.GetSparkle(c)
			if err != nil {
				return err
			}
			err = sparkle.Build(cmd.Context(), sparkleConfig)
			if err != nil {
				return err
			}

			return nil
		},
	}

	cmd.Flags().StringVarP(&configPath, "config", "c", "", "load configuration from a file")

	return cmd
}
