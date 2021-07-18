package main

import (
	"encoding/pem"
	"errors"
	"log"
	"os"
	"path/filepath"

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

	cmd.AddCommand(generateCmd(), publicCmd())

	return cmd
}

type keysOptions struct {
	path    string
	dsa     string
	ed25519 string
}

func generateCmd() *cobra.Command {
	opt := &keysOptions{}

	cmd := &cobra.Command{
		Use:     "generate",
		Short:   "Generate private keys",
		Long:    "Generate private keys for signing update packages.",
		Aliases: []string{"g"},
		Args:    cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			logErr := func(path string, fn func(string) error) error {
				err := fn(path)
				if err == os.ErrExist {
					log.Println("Key already exists:", path)
				} else if err == nil {
					log.Printf("Created DSA private key: %s", path)
				}
				return err
			}

			err := logErr(filepath.Join(opt.path, opt.dsa), createDSAPrivateKey)
			if err != nil {
				return err
			}

			return logErr(filepath.Join(opt.path, opt.ed25519), createEdPrivateKey)
		},
	}

	cmd.Flags().StringVarP(&opt.path, "path", "p", getDir(), "path to directory to create private keys in")
	cmd.Flags().StringVar(&opt.dsa, "dsa", "dsa.key", "file name of DSA private key")
	cmd.Flags().StringVar(&opt.ed25519, "ed25519", "ed25519.key", "file name of ed25519 private key")

	return cmd
}

func publicCmd() *cobra.Command {
	opt := &keysOptions{}

	cmd := &cobra.Command{
		Use:   "public <dsa|ed25519>",
		Short: "Output public key",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			var b []byte

			switch args[0] {
			case "dsa":
				dsaPath := filepath.Join(opt.path, opt.dsa)
				key, err := readKey(dsaPath, dsa.UnmarshalPrivateKey)
				if err != nil {
					return err
				}
				b, err = dsa.MarshalPublicKey(dsa.Public(key))
				if err != nil {
					return err
				}

			case "ed25519":
				edPath := filepath.Join(opt.path, opt.ed25519)
				key, err := readKey(edPath, ed25519.UnmarshalPrivateKey)
				if err != nil {
					return err
				}
				b, err = ed25519.MarshalPublicKey(ed25519.Public(key))
				if err != nil {
					return err
				}

			default:
				return errors.New("invalid argument '%s': should be 'dsa' or 'ed25519'")
			}

			return pem.Encode(os.Stdout, &pem.Block{Type: "PUBLIC KEY", Bytes: b})
		},
	}

	cmd.Flags().StringVarP(&opt.path, "path", "p", getDir(), "path to private keys")
	cmd.Flags().StringVar(&opt.dsa, "dsa", "dsa.key", "file name of DSA private key")
	cmd.Flags().StringVar(&opt.ed25519, "ed25519", "ed25519.key", "file name of ed25519 private key")

	return cmd
}

func createDSAPrivateKey(path string) error {
	if _, err := os.Stat(path); err != os.ErrNotExist {
		log.Println("Key already exists:", path)
		return nil
	}

	key, err := dsa.NewPrivateKey()
	if err != nil {
		return err
	}

	b, err := dsa.MarshalPrivateKey(key)
	if err != nil {
		return err
	}

	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()

	return pem.Encode(f, &pem.Block{Type: "PRIVATE KEY", Bytes: b})
}

func createEdPrivateKey(path string) error {
	if _, err := os.Stat(path); err != os.ErrNotExist {
		return os.ErrExist
	}

	key, err := ed25519.NewPrivateKey()
	if err != nil {
		return err
	}

	b, err := ed25519.MarshalPrivateKey(key)
	if err != nil {
		return err
	}

	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()

	return pem.Encode(f, &pem.Block{Type: "PRIVATE KEY", Bytes: b})
}
