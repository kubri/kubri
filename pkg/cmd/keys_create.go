package cmd

import (
	"errors"

	"github.com/abemedia/appcast/pkg/crypto/dsa"
	"github.com/abemedia/appcast/pkg/crypto/ed25519"
	"github.com/abemedia/appcast/pkg/crypto/pgp"
	"github.com/abemedia/appcast/pkg/crypto/rsa"
	"github.com/abemedia/appcast/pkg/secret"
	"github.com/spf13/cobra"
)

func keysCreateCmd() *cobra.Command {
	var name, email string

	cmd := &cobra.Command{
		Use:     "create",
		Short:   "Create private keys",
		Long:    "Create private keys for signing update packages. If keys already exist, this is a no-op.",
		Aliases: []string{"c"},
		Args:    cobra.NoArgs,
		RunE: func(*cobra.Command, []string) error {
			if err := createPrivateKey("dsa_key", dsa.NewPrivateKey, dsa.MarshalPrivateKey); err != nil {
				return err
			}
			if err := createPrivateKey("ed25519_key", ed25519.NewPrivateKey, ed25519.MarshalPrivateKey); err != nil {
				return err
			}
			if err := createPrivateKey("pgp_key", newPGPKey(name, email), pgp.MarshalPrivateKey); err != nil {
				return err
			}
			return createPrivateKey("rsa_key", rsa.NewPrivateKey, rsa.MarshalPrivateKey)
		},
	}

	cmd.Flags().StringVar(&name, "name", "", "your name for the PGP key")
	cmd.Flags().StringVar(&email, "email", "", "your email for the PGP key")

	return cmd
}

func newPGPKey(name, email string) func() (*pgp.PrivateKey, error) {
	return func() (*pgp.PrivateKey, error) {
		if name == "" && email == "" {
			return nil, errors.New("generating PGP key requires either name or email")
		}
		return pgp.NewPrivateKey(name, email)
	}
}

func createPrivateKey[PrivateKey any](
	name string,
	newKey func() (PrivateKey, error),
	marshal func(PrivateKey) ([]byte, error),
) error {
	if _, err := secret.Get(name); !errors.Is(err, secret.ErrKeyNotFound) {
		return err
	}
	key, err := newKey()
	if err != nil {
		return err
	}
	b, err := marshal(key)
	if err != nil {
		return err
	}
	return secret.Put(name, b)
}
