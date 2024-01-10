package cmd

import (
	"github.com/abemedia/appcast/pkg/pipe"
	"github.com/spf13/cobra"
)

func jsonschemaCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "jsonschema",
		Short: "Print config file jsonschema",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, _ []string) error {
			_, err := cmd.OutOrStdout().Write(pipe.Schema())
			return err
		},
	}

	return cmd
}
