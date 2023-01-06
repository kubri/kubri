package cmd

import (
	"os"

	"github.com/abemedia/appcast/pkg/secret"
	"github.com/spf13/cobra"
)

func keysImportCmd() *cobra.Command {
	var force bool

	cmd := &cobra.Command{
		Use:       "import (dsa|ed25519) <path>",
		Short:     "Import private keys",
		Long:      "Import private keys for signing update packages. If keys already exist, this is a no-op.",
		Aliases:   []string{"i"},
		Args:      cobra.ExactArgs(2),
		ValidArgs: []string{"dsa", "ed25519"},
		RunE: func(_ *cobra.Command, args []string) error {
			b, err := os.ReadFile(args[1])
			if err != nil {
				return err
			}

			key := args[0] + "_key"

			if force {
				_ = secret.Delete(key)
			}

			return secret.Put(key, b)
		},
	}

	cmd.Flags().BoolVarP(&force, "force", "f", false, "overwrite existing key")

	return cmd
}
