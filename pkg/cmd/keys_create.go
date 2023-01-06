package cmd

import (
	"encoding/pem"
	"errors"

	"github.com/abemedia/appcast/pkg/crypto/dsa"
	"github.com/abemedia/appcast/pkg/crypto/ed25519"
	"github.com/abemedia/appcast/pkg/secret"
	"github.com/spf13/cobra"
)

func keysCreateCmd() *cobra.Command {
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

				p := &pem.Block{Type: "PRIVATE KEY"}
				p.Bytes, err = dsa.MarshalPrivateKey(key)
				if err != nil {
					return err
				}

				if err = secret.Put("dsa_key", pem.EncodeToMemory(p)); err != nil {
					return err
				}
			}

			if _, err := secret.Get("ed25519_key"); errors.Is(err, secret.ErrKeyNotFound) {
				key, err := ed25519.NewPrivateKey()
				if err != nil {
					return err
				}

				p := &pem.Block{Type: "PRIVATE KEY"}
				p.Bytes, err = ed25519.MarshalPrivateKey(key)
				if err != nil {
					return err
				}

				if err = secret.Put("ed25519_key", pem.EncodeToMemory(p)); err != nil {
					return err
				}
			}

			return nil
		},
	}

	return cmd
}
