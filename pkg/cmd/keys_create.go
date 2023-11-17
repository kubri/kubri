package cmd

import (
	"errors"

	"github.com/abemedia/appcast/pkg/crypto/dsa"
	"github.com/abemedia/appcast/pkg/crypto/ed25519"
	"github.com/abemedia/appcast/pkg/crypto/pgp"
	"github.com/abemedia/appcast/pkg/secret"
	"github.com/spf13/cobra"
)

//nolint:funlen,gocognit
func keysCreateCmd() *cobra.Command {
	var name, email string

	cmd := &cobra.Command{
		Use:     "create",
		Short:   "Create private keys",
		Long:    "Create private keys for signing update packages. If keys already exist, this is a no-op.",
		Aliases: []string{"c"},
		Args:    cobra.NoArgs,
		RunE: func(*cobra.Command, []string) error {
			if _, err := secret.Get("dsa_key"); errors.Is(err, secret.ErrKeyNotFound) {
				key, err := dsa.NewPrivateKey()
				if err != nil {
					return err
				}

				b, err := dsa.MarshalPrivateKey(key)
				if err != nil {
					return err
				}

				if err = secret.Put("dsa_key", b); err != nil {
					return err
				}
			}

			if _, err := secret.Get("ed25519_key"); errors.Is(err, secret.ErrKeyNotFound) {
				key, err := ed25519.NewPrivateKey()
				if err != nil {
					return err
				}

				b, err := ed25519.MarshalPrivateKey(key)
				if err != nil {
					return err
				}

				if err = secret.Put("ed25519_key", b); err != nil {
					return err
				}
			}

			if _, err := secret.Get("pgp_key"); errors.Is(err, secret.ErrKeyNotFound) {
				if name == "" && email == "" {
					return errors.New("generating PGP key requires either name or email")
				}

				key, err := pgp.NewPrivateKey(name, email)
				if err != nil {
					return err
				}

				b, err := pgp.MarshalPrivateKey(key)
				if err != nil {
					return err
				}

				if err = secret.Put("pgp_key", b); err != nil {
					return err
				}
			}

			return nil
		},
	}

	cmd.Flags().StringVar(&name, "name", "", "your name for the PGP key")
	cmd.Flags().StringVar(&email, "email", "", "you email for the pgp key")

	return cmd
}
