package cmd

import (
	"encoding/pem"
	"os"

	"github.com/abemedia/appcast/pkg/crypto/dsa"
	"github.com/abemedia/appcast/pkg/crypto/ed25519"
	"github.com/abemedia/appcast/pkg/secret"
	"github.com/spf13/cobra"
)

func keysPublicCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:       "public (dsa|ed25519)",
		Short:     "Output public key",
		Aliases:   []string{"p"},
		Args:      cobra.ExactArgs(1),
		ValidArgs: []string{"dsa", "ed25519"},
		RunE: func(_ *cobra.Command, args []string) error {
			var b []byte
			switch args[0] {
			case "dsa":
				priv, err := secret.Get("dsa_key")
				if err != nil {
					return err
				}
				block, _ := pem.Decode(priv)
				key, err := dsa.UnmarshalPrivateKey(block.Bytes)
				if err != nil {
					return err
				}
				b, err = dsa.MarshalPublicKey(dsa.Public(key))
				if err != nil {
					return err
				}
			case "ed25519":
				priv, err := secret.Get("ed25519_key")
				if err != nil {
					return err
				}
				block, _ := pem.Decode(priv)
				key, err := ed25519.UnmarshalPrivateKey(block.Bytes)
				if err != nil {
					return err
				}
				b, err = ed25519.MarshalPublicKey(ed25519.Public(key))
				if err != nil {
					return err
				}
			}
			return pem.Encode(os.Stdout, &pem.Block{Type: "PUBLIC KEY", Bytes: b})
		},
	}

	return cmd
}
