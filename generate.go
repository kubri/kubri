package appcast

import (
	"bytes"
	"log"
	"mime"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"github.com/abemedia/appcast/source"
	"github.com/russross/blackfriday/v2"
	"golang.org/x/mod/semver"
)

func Generate(c *Config) (*Sparkle, error) {
	releases, err := c.Source.Releases()
	if err != nil {
		return nil, err
	}

	sort.Slice(releases, func(i, j int) bool { return releases[i].Date.After(releases[j].Date) })

	var items []SparkleItem
	for _, release := range releases {
		if !semver.IsValid(release.Version) {
			log.Println("skip invalid version", release.Version)
			continue
		}

		if release.Prerelease && !c.Prerelease {
			log.Println("skip prelease", release.Version)
			continue
		}

		item, err := releaseToSparkleItem(c, release)
		if err != nil {
			log.Println("warning:", err)
			continue
		}

		items = append(items, item...)
	}

	s := &Sparkle{
		Version:      "2.0",
		XMLNSSparkle: "http://www.andymatuschak.org/xml-namespaces/sparkle",
		XMLNSDC:      "http://purl.org/dc/elements/1.1/",
		Channels: []SparkleChannel{
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

func releaseToSparkleItem(c *Config, release *source.Release) ([]SparkleItem, error) {
	signatures, err := getSignatures(c, release)
	if err != nil {
		return nil, err
	}

	var description *CdataString
	if release.Description != "" {
		htmlDescription := blackfriday.Run([]byte(release.Description))
		description = &CdataString{string(bytes.TrimSpace(htmlDescription))}
	}

	items := make([]SparkleItem, 0, len(release.Assets))
	for _, asset := range release.Assets {
		os := detectOS(c, asset.Name)
		if os == Unknown {
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

		items = append(items, SparkleItem{
			Title:       release.Name,
			PubDate:     release.Date.UTC().Format(time.RFC1123),
			Description: description,
			Enclosure: SparkleEnclosure{
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

func detectOS(c *Config, url string) OS {
	if matchFallback(c.IsMacOS, isMacOS)(url) {
		return MacOS
	}
	if matchFallback(c.IsWindows64, isWindows64)(url) {
		return Windows64
	}
	if matchFallback(c.IsWindows32, isWindows32)(url) {
		return Windows32
	}

	return Unknown
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
