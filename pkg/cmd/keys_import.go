package cmd

import (
	"os"

	"github.com/abemedia/appcast/pkg/crypto/dsa"
	"github.com/abemedia/appcast/pkg/crypto/ed25519"
	"github.com/abemedia/appcast/pkg/crypto/pgp"
	"github.com/abemedia/appcast/pkg/crypto/rsa"
	"github.com/abemedia/appcast/pkg/secret"
	"github.com/spf13/cobra"
)

func keysImportCmd() *cobra.Command {
	var force bool

	cmd := &cobra.Command{
		Use:       "import (dsa|ed25519|pgp|rsa) <path>",
		Short:     "Import private keys",
		Long:      "Import private keys for signing update packages. If keys already exist, this is a no-op.",
		Aliases:   []string{"i"},
		Args:      cobra.ExactArgs(2),
		ValidArgs: []string{"dsa", "ed25519", "pgp", "rsa"},
		RunE: func(_ *cobra.Command, args []string) error {
			b, err := os.ReadFile(args[1])
			if err != nil {
				return err
			}

			switch args[0] {
			case "dsa":
				_, err = dsa.UnmarshalPrivateKey(b)
			case "ed25519":
				_, err = ed25519.UnmarshalPrivateKey(b)
			case "pgp":
				_, err = pgp.UnmarshalPrivateKey(b)
			case "rsa":
				_, err = rsa.UnmarshalPrivateKey(b)
			}
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
