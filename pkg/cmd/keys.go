package cmd

import (
	"fmt"
	"slices"
	"strings"

	"github.com/spf13/cobra"
)

func keysCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "keys",
		Short:   "Manage keys",
		Long:    "Manage keys for signing & verifying update packages.",
		Aliases: []string{"k"},
		Args:    cobra.NoArgs,
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			if len(args) > 0 && !slices.Contains(cmd.ValidArgs, args[0]) {
				var b strings.Builder
				for _, s := range cmd.ValidArgs {
					b.WriteString(fmt.Sprintf("\t%v\n", s))
				}
				return fmt.Errorf("invalid argument %q for %q\n\nDid you mean this?\n%s", args[0], cmd.CommandPath(), b.String())
			}
			return nil
		},
	}

	cmd.AddCommand(keysCreateCmd(), keysPublicCmd(), keysImportCmd(), keysExportCmd())

	return cmd
}
