package main

import (
	"path/filepath"

	"github.com/abemedia/appcast"
	"github.com/abemedia/appcast/pkg/crypto/dsa"
	"github.com/abemedia/appcast/pkg/crypto/ed25519"
	"github.com/spf13/cobra"
)

type signOptions struct {
	config  string
	source  string
	version string
	path    string
	dsaKey  string
	edKey   string
}

func signCmd() *cobra.Command {
	opt := &signOptions{}

	cmd := &cobra.Command{
		Use:     "sign",
		Short:   "Sign update packages",
		Long:    "Sign update packages and upload signatures to repository.",
		Aliases: []string{"s"},
		Args:    cobra.NoArgs,
		RunE: func(_ *cobra.Command, args []string) error {
			c, err := buildSign(opt)
			if err != nil {
				return err
			}

			return appcast.Sign(c)
		},
	}

	cmd.Flags().StringVarP(&opt.config, "config", "c", "", "load configuration from a file")
	cmd.Flags().StringVarP(&opt.source, "source", "s", "", "path or URL of files to sign")
	cmd.Flags().StringVarP(&opt.version, "version", "v", "",
		"version constraint to sign only specific releases e.g. 'v1.2.4', 'v1', '>= v1.1.0, < v2.1'")
	cmd.Flags().StringVarP(&opt.path, "path", "p", getDir(), "path to private keys")
	cmd.Flags().StringVar(&opt.dsaKey, "dsa", "dsa.key", "file name of DSA private key")
	cmd.Flags().StringVar(&opt.edKey, "ed25519", "ed25519.key", "file name of ed25519 private key")

	return cmd
}

func buildSign(opt *signOptions) (*appcast.SignOptions, error) {
	s := &appcast.SignOptions{}
	var err error

	if opt.config != "" {
		c, err := readConfig(opt.config)
		if err != nil {
			return nil, err
		}
		s.Source = c.Source
	}

	if opt.source != "" {
		s.Source, err = parseSource(opt.source)
		if err != nil {
			return nil, err
		}
	}

	s.DSAKey, err = readKey(filepath.Join(opt.path, opt.edKey), dsa.UnmarshalPrivateKey)
	if err != nil {
		return nil, err
	}

	s.Ed25519Key, err = readKey(filepath.Join(opt.path, opt.edKey), ed25519.UnmarshalPrivateKey)
	if err != nil {
		return nil, err
	}

	s.Version = opt.version

	return s, nil
}
