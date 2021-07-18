package appcast

import (
	"bytes"
	"log"
	"mime"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"github.com/abemedia/appcast/pkg/sparkle"
	"github.com/abemedia/appcast/source"
	"github.com/russross/blackfriday/v2"
)

// Feed generates an appcast feed.
func Feed(c *Config) (*sparkle.Feed, error) {
	releases, err := c.Source.ListReleases(nil)
	if err != nil {
		return nil, err
	}

	sort.Slice(releases, func(i, j int) bool { return releases[i].Date.After(releases[j].Date) })

	var items []sparkle.Item
	for _, release := range releases {
		if release.Prerelease && !c.Prerelease {
			log.Println("skipping prelease:", release.Version)
			continue
		}

		item, err := releaseToSparkleItem(c, release)
		if err != nil {
			log.Println("warning:", err)
			continue
		}

		items = append(items, item...)
	}

	s := &sparkle.Feed{
		Version:      "2.0",
		XMLNSSparkle: "http://www.andymatuschak.org/xml-namespaces/sparkle",
		XMLNSDC:      "http://purl.org/dc/elements/1.1/",
		Channels: []sparkle.Channel{
			{
				Title:       c.Title,
				Link:        c.URL,
				Description: c.Description,
				Items:       items,
			},
		},
	}

	return s, nil
}

func releaseToSparkleItem(c *Config, release *source.Release) ([]sparkle.Item, error) {
	signatures, err := getSignatures(c, release)
	if err != nil {
		return nil, err
	}

	var description *sparkle.CdataString
	if release.Description != "" {
		htmlDescription := blackfriday.Run([]byte(release.Description))
		description = &sparkle.CdataString{Value: string(bytes.TrimSpace(htmlDescription))}
	}

	items := make([]sparkle.Item, 0, len(release.Assets))
	for _, asset := range release.Assets {
		os := detectOS(asset.Name)
		if os == Unknown {
			log.Printf("skipping %s: unsupported file extension\n", asset.Name)
			continue
		}

		opt, err := getSettings(c.Settings, release.Version, os)
		if err != nil {
			return nil, err
		}

		url := asset.URL
		if c.RewriteURL != nil {
			url = c.RewriteURL(url)
		}

		mimeType := mime.TypeByExtension(filepath.Ext(asset.Name))
		if mimeType == "" {
			mimeType = "application/octet-stream"
		}

		items = append(items, sparkle.Item{
			Title:       release.Name,
			PubDate:     release.Date.UTC().Format(time.RFC1123),
			Description: description,
			Enclosure: sparkle.Enclosure{
				Version:              strings.TrimPrefix(release.Version, "v"),
				URL:                  url,
				InstallerArguments:   opt.InstallerArguments,
				MinimumSystemVersion: opt.MinimumSystemVersion,
				Type:                 mimeType,
				OS:                   os.String(),
				Length:               asset.Size,
				DSASignature:         signatures.Get(asset.Name, "dsa"),
				EDSignature:          signatures.Get(asset.Name, "ed25519"),
			},
		})
	}

	return items, nil
}

func getSignatures(c *Config, release *source.Release) (signatures, error) {
	sig := signatures{}

	s := getAsset(release.Assets, "signatures.txt")
	if s == nil {
		return sig, nil
	}

	b, err := c.Source.DownloadAsset(release.Version, s.Name)
	if err != nil {
		return nil, err
	}

	if err = sig.UnmarshalText(b); err != nil {
		return nil, err
	}

	return sig, nil
}

func getAsset(assets []*source.Asset, name string) *source.Asset {
	for _, asset := range assets {
		if strings.HasSuffix(asset.URL, name) {
			return asset
		}
	}
	return nil
}
