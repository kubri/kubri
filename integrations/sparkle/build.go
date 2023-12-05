package sparkle

import (
	"context"
	"crypto/sha1"
	"encoding/base64"
	"encoding/xml"
	"path"
	"strings"
	"time"
	"unsafe"

	"github.com/abemedia/appcast/pkg/crypto/dsa"
	"github.com/abemedia/appcast/pkg/crypto/ed25519"
	"github.com/abemedia/appcast/source"
	"github.com/go-xmlfmt/xmlfmt"
	"github.com/russross/blackfriday/v2"
)

func Build(ctx context.Context, c *Config) error {
	items := read(ctx, c)

	version := c.Version
	if v := getVersionConstraint(items); v != "" {
		version += "," + v
	}

	releases, err := c.Source.ListReleases(ctx, &source.ListOptions{
		Version:    version,
		Prerelease: c.Prerelease,
	})
	if err == source.ErrNoReleaseFound {
		return nil
	}
	if err != nil {
		return err
	}

	i, err := getItems(ctx, c, releases)
	if err != nil {
		return err
	}
	items = append(i, items...)

	link, err := c.Target.URL(ctx, c.FileName)
	if err != nil {
		return err
	}

	rss := &RSS{Channels: []*Channel{{
		Title:       c.Title,
		Description: c.Description,
		Link:        link,
		Items:       items,
	}}}

	return write(ctx, c, rss)
}

func read(ctx context.Context, c *Config) []*Item {
	r, err := c.Target.NewReader(ctx, c.FileName)
	if err != nil {
		return nil
	}
	defer r.Close()
	var rss RSS
	if err = xml.NewDecoder(r).Decode(&rss); err != nil || len(rss.Channels) == 0 {
		return nil
	}
	return rss.Channels[0].Items
}

func getVersionConstraint(items []*Item) string {
	if len(items) == 0 {
		return ""
	}

	v := make([]byte, 0, len(items)*len("!=0.0.0,"))
	for _, item := range items {
		v = append(v, '!', '=')
		v = append(v, item.Version...)
		v = append(v, ',')
	}

	return unsafe.String(unsafe.SliceData(v), len(v)-1)
}

func getItems(ctx context.Context, c *Config, releases []*source.Release) ([]*Item, error) {
	var items []*Item
	for _, release := range releases {
		item, err := getReleaseItems(ctx, c, release)
		if err != nil {
			return nil, err
		}
		items = append(items, item...)
	}
	return items, nil
}

//nolint:funlen
func getReleaseItems(ctx context.Context, c *Config, release *source.Release) ([]*Item, error) {
	var description *CdataString
	if release.Description != "" {
		desc := string(blackfriday.Run([]byte(release.Description)))
		desc = strings.TrimSpace(xmlfmt.FormatXML(desc, "\t\t\t\t", "\t"))
		description = &CdataString{Value: "\n\t\t\t\t" + desc + "\n\t\t\t"}
	}

	items := make([]*Item, 0, len(release.Assets))
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

		var b []byte
		if (os == MacOS && c.Ed25519Key != nil) || c.DSAKey != nil || c.UploadPackages {
			b, err = c.Source.DownloadAsset(ctx, release.Version, asset.Name)
			if err != nil {
				return nil, err
			}
		}

		edSig, dsaSig, err := signAsset(c, os, b)
		if err != nil {
			return nil, err
		}

		url := asset.URL
		if c.UploadPackages {
			url, err = uploadAsset(ctx, c, release.Version+"/"+asset.Name, b)
			if err != nil {
				return nil, err
			}
		}

		version := strings.TrimPrefix(release.Version, "v")

		items = append(items, &Item{
			Title:                             release.Name,
			PubDate:                           release.Date.UTC().Format(time.RFC1123),
			Description:                       description,
			Version:                           version,
			CriticalUpdate:                    getCriticalUpdate(opt.CriticalUpdateBelowVersion),
			Tags:                              getTags(opt.CriticalUpdate),
			IgnoreSkippedUpgradesBelowVersion: opt.IgnoreSkippedUpgradesBelowVersion,
			MinimumAutoupdateVersion:          opt.MinimumAutoupdateVersion,
			Enclosure: &Enclosure{
				Version:              version,
				URL:                  url,
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
func signAsset(c *Config, os OS, b []byte) (edSig, dsaSig string, err error) {
	if os == MacOS && c.Ed25519Key != nil {
		sig, err := ed25519.Sign(c.Ed25519Key, b)
		if err != nil {
			return "", "", err
		}
		edSig = base64.RawStdEncoding.EncodeToString(sig)
	} else if c.DSAKey != nil {
		sum := sha1.Sum(b)
		sig, err := dsa.Sign(c.DSAKey, sum[:])
		if err != nil {
			return "", "", err
		}
		dsaSig = base64.RawStdEncoding.EncodeToString(sig)
	}
	return
}

func uploadAsset(ctx context.Context, c *Config, name string, b []byte) (string, error) {
	w, err := c.Target.NewWriter(ctx, name)
	if err != nil {
		return "", err
	}
	if _, err := w.Write(b); err != nil {
		return "", err
	}
	if err = w.Close(); err != nil {
		return "", err
	}
	return c.Target.URL(ctx, name)
}

func getCriticalUpdate(version string) *CriticalUpdate {
	if version != "" {
		return &CriticalUpdate{Version: version}
	}
	return nil
}

func getTags(critical bool) *Tags {
	if critical {
		return &Tags{CriticalUpdate: true}
	}
	return nil
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

//nolint:gochecknoglobals
var replacer = strings.NewReplacer("></sparkle:criticalUpdate>", " />", "></enclosure>", " />")

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

//nolint:gochecknoinits // Don't use carriage return on windows.
func init() {
	xmlfmt.NL = "\n"
}
