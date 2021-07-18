package main

import (
	"encoding/pem"
	"errors"
	"os"

	"github.com/abemedia/appcast/pkg/crypto/dsa"
	"github.com/abemedia/appcast/pkg/crypto/ed25519"
	"github.com/spf13/cobra"
)

func keysCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "keys",
		Short:   "Generate DSA & ed25519 keys",
		Long:    "Generate DSA & ed25519 keys for signing & verifying update packages.",
		Aliases: []string{"k"},
		Args:    cobra.NoArgs,
	}

	cmd.AddCommand(privateCmd(), publicCmd())

	return cmd
}

func privateCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "private <dsa|ed>",
		Short: "Generate a private key",
		Long:  "Generate a private key for signing update packages.",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			var b []byte

			switch args[0] {
			case "dsa":
				key, err := dsa.NewPrivateKey()
				if err != nil {
					return err
				}
				b, err = dsa.MarshalPrivateKey(key)
				if err != nil {
					return err
				}

			case "ed":
				key, err := ed25519.NewPrivateKey()
				if err != nil {
					return err
				}
				b, err = ed25519.MarshalPrivateKey(key)
				if err != nil {
					return err
				}

			default:
				return errors.New("invalid argument '%s': should be 'dsa' or 'ed'")
			}

			return pem.Encode(os.Stdout, &pem.Block{Type: "PRIVATE KEY", Bytes: b})
		},
	}

	return cmd
}

func publicCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "public <dsa|ed> <path>",
		Short: "Generate public key from private key",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			b, err := os.ReadFile(args[1])
			if err != nil {
				return err
			}

			block, _ := pem.Decode(b)

			switch args[0] {
			case "dsa":
				key, err := dsa.UnmarshalPrivateKey(block.Bytes)
				if err != nil {
					return err
				}
				b, err = dsa.MarshalPublicKey(dsa.NewPublicKey(key))
				if err != nil {
					return err
				}

			case "ed":
				key, err := ed25519.UnmarshalPrivateKey(block.Bytes)
				if err != nil {
					return err
				}
				b, err = ed25519.MarshalPublicKey(ed25519.NewPublicKey(key))
				if err != nil {
					return err
				}

			default:
				return errors.New("invalid argument '%s': should be 'dsa' or 'ed'")
			}

			return pem.Encode(os.Stdout, &pem.Block{Type: "PUBLIC KEY", Bytes: b})
		},
	}

	return cmd
}
