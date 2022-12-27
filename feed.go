package appcast

import (
	"bytes"
	"context"
	"log"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"github.com/abemedia/appcast/pkg/os"
	"github.com/abemedia/appcast/pkg/sparkle"
	"github.com/abemedia/appcast/source"
	"github.com/russross/blackfriday/v2"
)

// Feed generates an appcast feed.
func Feed(ctx context.Context, c *Config) (*sparkle.Feed, error) {
	releases, err := c.Source.ListReleases(ctx, nil)
	if err != nil {
		return nil, err
	}

	sort.Slice(releases, func(i, j int) bool { return releases[i].Date.After(releases[j].Date) })

	var items []sparkle.Item
	for _, release := range releases {
		if release.Prerelease && !c.Prerelease {
			log.Println("Skipping prelease:", release.Version)
			continue
		}

		item, err := releaseToSparkleItem(ctx, c, release)
		if err != nil {
			log.Printf("Skipping %s: %s", release.Version, err)
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

func releaseToSparkleItem(ctx context.Context, c *Config, release *source.Release) ([]sparkle.Item, error) {
	signatures, err := getSignatures(ctx, c, release)
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
		o := os.Detect(asset.Name)
		if o == os.Unknown {
			log.Printf("Skipping asset %s (%s): unsupported file extension\n", asset.Name, release.Version)
			continue
		}

		opt, err := getSettings(c.Settings, release.Version, o)
		if err != nil {
			return nil, err
		}

		url := asset.URL
		if c.RewriteURL != nil {
			url = c.RewriteURL(url)
		}

		version := strings.TrimPrefix(release.Version, "v")

		items = append(items, sparkle.Item{
			Title:                             release.Name,
			PubDate:                           release.Date.UTC().Format(time.RFC1123),
			Description:                       description,
			Version:                           version,
			CriticalUpdate:                    getCriticalUpdate(opt.CriticalUpdateBelowVersion),
			Tags:                              getTags(opt.CriticalUpdate),
			IgnoreSkippedUpgradesBelowVersion: opt.IgnoreSkippedUpgradesBelowVersion,
			MinimumAutoupdateVersion:          opt.MinimumAutoupdateVersion,
			Enclosure: sparkle.Enclosure{
				Version:              version,
				URL:                  url,
				InstallerArguments:   opt.InstallerArguments,
				MinimumSystemVersion: opt.MinimumSystemVersion,
				Type:                 getFileType(asset.Name),
				OS:                   o.String(),
				Length:               asset.Size,
				DSASignature:         signatures.Get(asset.Name, "dsa"),
				EDSignature:          signatures.Get(asset.Name, "ed25519"),
			},
		})
	}

	return items, nil
}

func getFileType(s string) string {
	ext := filepath.Ext(s)
	switch ext {
	default:
		return "application/octet-stream"
	case ".pkg", ".mpkg":
		return "application/vnd.apple.installer+xml"
	case ".dmg":
		return "application/x-apple-diskimage"
	case ".msi":
		return "application/x-msi"
	case ".exe":
		return "application/vnd.microsoft.portable-executable"

	// See https://learn.microsoft.com/en-us/windows/msix/app-installer/web-install-iis#step-7---configure-the-web-app-for-app-package-mime-types
	case ".msix", ".msixbundle", ".appx", ".appxbundle", ".appinstaller":
		return "application/" + ext[1:]
	}
}

func getTags(critical bool) *sparkle.Tags {
	if critical {
		return &sparkle.Tags{CriticalUpdate: true}
	}
	return nil
}

func getCriticalUpdate(version string) *sparkle.CriticalUpdate {
	if version != "" {
		return &sparkle.CriticalUpdate{Version: version}
	}
	return nil
}

func getSignatures(ctx context.Context, c *Config, release *source.Release) (signatures, error) {
	sig := signatures{}

	s := getAsset(release.Assets, "signatures.txt")
	if s == nil {
		return sig, nil
	}

	b, err := c.Source.DownloadAsset(ctx, release.Version, s.Name)
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
