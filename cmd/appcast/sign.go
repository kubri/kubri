package main

import (
	"os"

	"github.com/abemedia/appcast"
	"github.com/abemedia/appcast/source"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
)

type signOptions struct {
	config string
	source string
	dsaKey string
	edKey  string
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
	cmd.Flags().StringVarP(&opt.source, "source", "s", "", "source of files to sign")
	cmd.Flags().StringVar(&opt.dsaKey, "dsa", "", "path to DSA private key")
	cmd.Flags().StringVar(&opt.edKey, "ed", "", "path to ed25519 private key")

	return cmd
}

func buildSign(opt *signOptions) (*appcast.Config, error) {
	c := &appcast.Config{}

	if opt.config != "" {
		b, err := os.ReadFile(opt.config)
		if err != nil {
			return nil, err
		}
		if err = yaml.Unmarshal(b, c); err != nil {
			return nil, err
		}
	}

	if opt.source != "" {
		if _, err := os.Stat(opt.source); err == nil {
			opt.source = "local://" + opt.source
		}
		c.Source = &source.Source{}
		err := c.Source.UnmarshalText([]byte(opt.source))
		if err != nil {
			return nil, err
		}
	}

	if opt.dsaKey != "" {
		c.DSAKey = opt.dsaKey
	}

	if opt.edKey != "" {
		c.EdKey = opt.edKey
	}

	return c, nil
}
