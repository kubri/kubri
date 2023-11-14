package cmd

import (
	"os"

	"github.com/abemedia/appcast/pkg/crypto/dsa"
	"github.com/abemedia/appcast/pkg/crypto/ed25519"
	"github.com/abemedia/appcast/pkg/crypto/pgp"
	"github.com/abemedia/appcast/pkg/secret"
	"github.com/spf13/cobra"
)

func keysPublicCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:       "public (dsa|ed25519|pgp)",
		Short:     "Output public key",
		Aliases:   []string{"p"},
		Args:      cobra.ExactArgs(1),
		ValidArgs: []string{"dsa", "ed25519", "pgp"},
		RunE: func(_ *cobra.Command, args []string) error {
			var pub []byte
			switch args[0] {
			case "dsa":
				priv, err := secret.Get("dsa_key")
				if err != nil {
					return err
				}
				key, err := dsa.UnmarshalPrivateKey(priv)
				if err != nil {
					return err
				}
				pub, err = dsa.MarshalPublicKey(dsa.Public(key))
				if err != nil {
					return err
				}
			case "ed25519":
				priv, err := secret.Get("ed25519_key")
				if err != nil {
					return err
				}
				key, err := ed25519.UnmarshalPrivateKey(priv)
				if err != nil {
					return err
				}
				pub, err = ed25519.MarshalPublicKey(ed25519.Public(key))
				if err != nil {
					return err
				}
			case "pgp":
				priv, err := secret.Get("pgp_key")
				if err != nil {
					return err
				}
				key, err := pgp.UnmarshalPrivateKey(priv)
				if err != nil {
					return err
				}
				pub, err = pgp.MarshalPublicKey(pgp.Public(key))
				if err != nil {
					return err
				}
			}

			_, err := os.Stdout.Write(pub)
			return err
		},
	}

	return cmd
}
