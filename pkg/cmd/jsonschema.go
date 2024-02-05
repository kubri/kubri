package cmd

import (
	"github.com/abemedia/appcast/pkg/config"
	"github.com/spf13/cobra"
)

func jsonschemaCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "jsonschema",
		Short: "Print config file jsonschema",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, _ []string) error {
			_, err := cmd.OutOrStdout().Write(config.Schema())
			return err
		},
	}

	return cmd
}
