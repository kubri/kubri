package sparkle

import (
	"context"
	"crypto/sha1"
	"encoding/base64"
	"encoding/xml"
	"path"
	"strings"
	"time"

	"github.com/abemedia/appcast/pkg/crypto/dsa"
	"github.com/abemedia/appcast/pkg/crypto/ed25519"
	"github.com/abemedia/appcast/source"
	"github.com/go-xmlfmt/xmlfmt"
	"github.com/russross/blackfriday/v2"
)

func Build(ctx context.Context, c *Config) error {
	releases, err := c.Source.ListReleases(ctx, &source.ListOptions{
		Version:    c.Version,
		Prerelease: c.Prerelease,
	})
	if err != nil {
		return err
	}

	cached := map[string][]Item{}
	if r, err := read(ctx, c); err == nil {
		for _, item := range r.Channels[0].Items {
			cached[item.Version] = append(cached[item.Version], item)
		}
	}

	var items []Item
	for _, release := range releases {
		// Load saved RSS.
		if item, ok := cached[release.Version[1:]]; ok {
			items = append(items, item...)
			continue
		}

		item, err := createReleaseItems(ctx, c, release)
		if err != nil {
			return err
		}

		items = append(items, item...)
	}

	return write(ctx, c, newRSS(c.Title, c.Description, c.URL, items))
}

//nolint:gochecknoglobals
var replacer = strings.NewReplacer("></sparkle:criticalUpdate>", " />", "></enclosure>", " />")

func read(ctx context.Context, c *Config) (*RSS, error) {
	r, err := c.Target.NewReader(ctx, c.FileName)
	if err != nil {
		return nil, err
	}
	defer r.Close()
	var rss RSS
	if err = xml.NewDecoder(r).Decode(&rss); err != nil {
		return nil, err
	}
	return &rss, nil
}

func write(ctx context.Context, c *Config, rss *RSS) error {
	w, err := c.Target.NewWriter(ctx, c.FileName)
	if err != nil {
		return err
	}

	_, err = w.Write([]byte(xml.Header))
	if err != nil {
		return err
	}

	b, err := xml.MarshalIndent(rss, "", "\t")
	if err != nil {
		return err
	}

	_, err = replacer.WriteString(w, string(b))
	if err != nil {
		return err
	}

	return w.Close()
}

func newRSS(title, description, url string, items []Item) *RSS {
	return &RSS{Channels: []Channel{{Title: title, Description: description, Link: url, Items: items}}}
}

func createReleaseItems(ctx context.Context, c *Config, release *source.Release) ([]Item, error) {
	var description *CdataString
	if release.Description != "" {
		desc := string(blackfriday.Run([]byte(release.Description)))
		desc = strings.TrimSpace(xmlfmt.FormatXML(desc, "\t\t\t\t", "\t"))
		description = &CdataString{Value: "\n\t\t\t\t" + desc + "\n\t\t\t"}
	}

	items := make([]Item, 0, len(release.Assets))
	for _, asset := range release.Assets {
		detect := c.DetectOS
		if detect == nil {
			detect = DetectOS
		}
		os := detect(asset.Name)
		if os == Unknown {
			continue
		}

		opt, err := getSettings(c.Settings, release.Version, os)
		if err != nil {
			return nil, err
		}

		edSig, dsaSig, err := signAsset(ctx, c, release.Version, os, asset)
		if err != nil {
			return nil, err
		}

		version := strings.TrimPrefix(release.Version, "v")

		items = append(items, Item{
			Title:                             release.Name,
			PubDate:                           release.Date.UTC().Format(time.RFC1123),
			Description:                       description,
			Version:                           version,
			CriticalUpdate:                    getCriticalUpdate(opt.CriticalUpdateBelowVersion),
			Tags:                              getTags(opt.CriticalUpdate),
			IgnoreSkippedUpgradesBelowVersion: opt.IgnoreSkippedUpgradesBelowVersion,
			MinimumAutoupdateVersion:          opt.MinimumAutoupdateVersion,
			Enclosure: Enclosure{
				Version:              version,
				URL:                  asset.URL,
				InstallerArguments:   opt.InstallerArguments,
				MinimumSystemVersion: opt.MinimumSystemVersion,
				Type:                 getFileType(asset.Name),
				OS:                   os.String(),
				Length:               asset.Size,
				DSASignature:         dsaSig,
				EDSignature:          edSig,
			},
		})
	}

	return items, nil
}

//nolint:nonamedreturns
func signAsset(ctx context.Context, c *Config, v string, os OS, a *source.Asset) (edSig, dsaSig string, err error) {
	if os == MacOS && c.Ed25519Key != nil {
		b, err := c.Source.DownloadAsset(ctx, v, a.Name)
		if err != nil {
			return "", "", err
		}
		sig := ed25519.Sign(c.Ed25519Key, b)
		edSig = base64.RawStdEncoding.EncodeToString(sig)
	} else if c.DSAKey != nil {
		b, err := c.Source.DownloadAsset(ctx, v, a.Name)
		if err != nil {
			return "", "", err
		}
		sum := sha1.Sum(b)
		sum = sha1.Sum(sum[:])
		sig, err := dsa.Sign(c.DSAKey, sum[:])
		if err != nil {
			return "", "", err
		}
		dsaSig = base64.RawStdEncoding.EncodeToString(sig)
	}
	return
}

func getFileType(s string) string {
	ext := path.Ext(s)
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

func getTags(critical bool) *Tags {
	if critical {
		return &Tags{CriticalUpdate: true}
	}
	return nil
}

func getCriticalUpdate(version string) *CriticalUpdate {
	if version != "" {
		return &CriticalUpdate{Version: version}
	}
	return nil
}
