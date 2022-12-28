package main

import (
	"encoding/xml"
	"io"
	"os"
	"path"
	"strings"

	"github.com/abemedia/appcast"
	"github.com/abemedia/appcast/source"
	"github.com/spf13/cobra"
	"golang.org/x/mod/semver"
	"gopkg.in/yaml.v3"
)

type feedOptions struct {
	config     string
	source     string
	target     string
	version    string
	prerelease bool

	title       string
	description string
	url         string
}

func feedCmd() *cobra.Command {
	opt := &feedOptions{}

	cmd := &cobra.Command{
		Use:     "feed",
		Short:   "Generates appcast feed XML",
		Aliases: []string{"f"},
		Args:    cobra.NoArgs,
		RunE: func(_ *cobra.Command, args []string) error {
			c, err := buildFeed(opt)
			if err != nil {
				return err
			}

			target, err := buildTarget(opt, c)
			if err != nil {
				return err
			}

			feed, err := appcast.Feed(appContext(), c)
			if err != nil {
				return err
			}

			_, err = target.Write([]byte(xml.Header))
			if err != nil {
				return err
			}

			enc := xml.NewEncoder(target)
			enc.Indent("", "  ")

			err = enc.Encode(feed)
			if err != nil {
				return err
			}

			return target.Close()
		},
	}

	cmd.Flags().StringVarP(&opt.config, "config", "c", "", "load configuration from a file")
	cmd.Flags().StringVarP(&opt.source, "source", "s", "", "source of releases")
	cmd.Flags().StringVarP(&opt.target, "target", "t", "", "where to store appcast XML")
	cmd.Flags().StringVarP(&opt.version, "version", "v", "", "version constraint to include only specific releases")
	cmd.Flags().BoolVar(&opt.prerelease, "prerelease", false, "include prereleases")

	cmd.Flags().StringVar(&opt.title, "title", "", "title of the appcast")
	cmd.Flags().StringVar(&opt.description, "description", "", "description of the appcast")
	cmd.Flags().StringVar(&opt.url, "url", "", "url of the appcast")

	return cmd
}

func buildTarget(opt *feedOptions, c *appcast.Config) (io.WriteCloser, error) {
	if opt.target == "" {
		return os.Stdout, nil
	}

	var version string
	if strings.HasPrefix(opt.version, semver.Canonical(opt.version)) {
		version = opt.version
	}

	filename := "appcast.xml"

	if opt.target == "source" {
		return source.NewWriter(c.Source, version, filename), nil
	}

	target := opt.target
	if strings.HasSuffix(opt.target, ".xml") {
		target, filename = path.Split(opt.target)
	}

	src, err := parseSource(target)
	if err != nil {
		return nil, err
	}

	return source.NewWriter(src, version, filename), nil
}

func buildFeed(opt *feedOptions) (*appcast.Config, error) {
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
		src, err := parseSource(opt.source)
		if err != nil {
			return nil, err
		}
		c.Source = src
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
