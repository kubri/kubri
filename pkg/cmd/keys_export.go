package cmd

import (
	"github.com/spf13/cobra"

	"github.com/kubri/kubri/pkg/crypto/ed25519"
	"github.com/kubri/kubri/pkg/secret"
)

func keysExportCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:       "export (dsa|ed25519|pgp|rsa)",
		Short:     "Export private keys",
		Long:      "Export private keys for signing update packages.",
		Aliases:   []string{"e"},
		Args:      cobra.ExactArgs(1),
		ValidArgs: []string{"dsa", "ed25519", "pgp", "rsa"},
		RunE: func(cmd *cobra.Command, args []string) error {
			out, err := secret.Get(args[0] + "_key")
			if err != nil {
				return err
			}

			// Convert ed25519 key to PEM format for better compatibility.
			if args[0] == "ed25519" { //nolint:goconst
				key, err := ed25519.UnmarshalPrivateKey(out)
				if err != nil {
					return err
				}
				out, err = ed25519.MarshalPrivateKeyPEM(key)
				if err != nil {
					return err
				}
			}

			_, err = cmd.OutOrStdout().Write(out)
			return err
		},
	}

	return cmd
}
