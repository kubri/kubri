package main

import (
	"encoding/xml"
	"os"

	"github.com/abemedia/appcast"
	"github.com/abemedia/appcast/source"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
)

type generateOptions struct {
	config     string
	source     string
	prerelease bool

	title       string
	description string
	url         string
}

func generateCmd() *cobra.Command {
	opt := &generateOptions{}

	cmd := &cobra.Command{
		Use:     "generate",
		Short:   "Generates appcast XML",
		Aliases: []string{"g"},
		Args:    cobra.ExactArgs(1),
		RunE: func(_ *cobra.Command, args []string) error {
			c, err := buildGenerate(opt)
			if err != nil {
				return err
			}

			sparkle, err := appcast.Generate(c)
			if err != nil {
				return err
			}

			enc := xml.NewEncoder(os.Stdout)
			enc.Indent("", "  ")

			return enc.Encode(sparkle)
		},
	}

	cmd.Flags().StringVarP(&opt.config, "config", "c", "appcast.yml", "load configuration from a file")
	cmd.Flags().StringVarP(&opt.source, "source", "s", "", "source of files to sign")
	cmd.Flags().BoolVar(&opt.prerelease, "prerelease", false, "include prereleases")

	cmd.Flags().StringVar(&opt.title, "title", "", "title of the appcast")
	cmd.Flags().StringVar(&opt.description, "description", "", "description of the appcast")
	cmd.Flags().StringVar(&opt.url, "url", "", "url of the appcast")

	return cmd
}

func buildGenerate(opt *generateOptions) (*appcast.Config, error) {
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
			opt.source = "file://" + opt.source
		}
		c.Source = &source.Source{}
		err := c.Source.UnmarshalText([]byte(opt.source))
		if err != nil {
			return nil, err
		}
	}

	if opt.prerelease {
		c.Prerelease = true
	}

	if opt.title != "" {
		c.Title = opt.title
	}

	if opt.description != "" {
		c.Description = opt.description
	}

	if opt.url != "" {
		c.URL = opt.url
	}

	return c, nil
}
