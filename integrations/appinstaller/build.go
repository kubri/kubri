package appinstaller

import (
	"archive/zip"
	"bytes"
	"context"
	"encoding/xml"
	"errors"
	"path"
	"strings"

	"github.com/abemedia/appcast/source"
	"github.com/abemedia/appcast/target"
)

type Config struct {
	Source                    *source.Source
	Target                    target.Target
	Version                   string
	Prerelease                bool
	HoursBetweenUpdateChecks  int
	UpdateBlocksActivation    bool
	ShowPrompt                bool
	AutomaticBackgroundTask   bool
	ForceUpdateFromAnyVersion bool
}

var ErrNotValid = errors.New("not a valid bundle")

func Build(ctx context.Context, c *Config) error {
	releases, err := c.Source.ListReleases(ctx, &source.ListOptions{
		Version:    c.Version,
		Prerelease: c.Prerelease,
	})
	if err != nil {
		return err
	}

	var r *source.Release
	for _, r = range releases {
		var ok bool
		for _, asset := range r.Assets {
			switch path.Ext(asset.Name) {
			case ".msix", ".appx", ".msixbundle", "appxbundle":
				res, err := build(ctx, c, r.Version, asset)
				if err != nil {
					return err
				}
				if err = write(ctx, c, res); err != nil {
					return err
				}
				ok = true
			}
		}
		if ok {
			break // Only build for latest release containing supported files.
		}
	}

	return nil
}

func build(ctx context.Context, c *Config, version string, asset *source.Asset) (*AppInstaller, error) {
	b, err := c.Source.DownloadAsset(ctx, version, asset.Name)
	if err != nil {
		return nil, err
	}

	p, err := getPackage(b)
	if err != nil {
		return nil, err
	}
	p.URI = asset.URL

	res := newXML(c)
	res.Version = p.Version

	var name string
	if strings.HasSuffix(asset.Name, "bundle") {
		res.MainBundle = p
		name = p.Name + ".appinstaller"
	} else {
		res.MainPackage = p
		name = p.Name + "-" + p.ProcessorArchitecture + ".appinstaller"
	}

	res.URI, err = c.Target.URL(ctx, name)
	if err != nil {
		return nil, err
	}

	return res, nil
}

func newXML(c *Config) *AppInstaller {
	appInstaller := &AppInstaller{}

	// Get minimum required namespace.
	switch {
	case c.ForceUpdateFromAnyVersion, c.ShowPrompt, c.UpdateBlocksActivation:
		appInstaller.XMLName.Space = "http://schemas.microsoft.com/appx/appinstaller/2018"
	case c.AutomaticBackgroundTask, c.HoursBetweenUpdateChecks > 0:
		appInstaller.XMLName.Space = "http://schemas.microsoft.com/appx/appinstaller/2017/2"
	default:
		appInstaller.XMLName.Space = "http://schemas.microsoft.com/appx/appinstaller/2017"
		return appInstaller
	}

	if c.AutomaticBackgroundTask || c.ForceUpdateFromAnyVersion {
		appInstaller.UpdateSettings = &UpdateSettings{
			AutomaticBackgroundTask:   Bool(c.AutomaticBackgroundTask),
			ForceUpdateFromAnyVersion: c.ForceUpdateFromAnyVersion,
		}
	}
	if c.HoursBetweenUpdateChecks > 0 || c.ShowPrompt || c.UpdateBlocksActivation {
		if appInstaller.UpdateSettings == nil {
			appInstaller.UpdateSettings = &UpdateSettings{}
		}
		appInstaller.UpdateSettings.OnLaunch = &OnLaunch{
			HoursBetweenUpdateChecks: c.HoursBetweenUpdateChecks,
			UpdateBlocksActivation:   c.UpdateBlocksActivation,
			ShowPrompt:               c.ShowPrompt,
		}
	}

	return appInstaller
}

func getPackage(b []byte) (*Package, error) {
	r, err := zip.NewReader(bytes.NewReader(b), int64(len(b)))
	if err != nil {
		return nil, err
	}
	for _, f := range r.File {
		if f.Name == "AppxManifest.xml" || f.Name == "AppxBundleManifest.xml" {
			rc, err := f.Open()
			if err != nil {
				return nil, err
			}
			defer rc.Close()
			var manifest struct{ Identity *Package }
			if err = xml.NewDecoder(rc).Decode(&manifest); err != nil {
				return nil, err
			}
			return manifest.Identity, nil
		}
	}
	return nil, ErrNotValid
}

//nolint:gochecknoglobals
var replacer = strings.NewReplacer("></Package>", " />", "></Bundle>", " />",
	"></MainPackage>", " />", "></MainBundle>", " />", "></OnLaunch>", " />",
	"></AutomaticBackgroundTask>", " />", "></ForceUpdateFromAnyVersion>", " />")

func write(ctx context.Context, c *Config, app *AppInstaller) error {
	w, err := c.Target.NewWriter(ctx, path.Base(app.URI))
	if err != nil {
		return err
	}

	_, err = w.Write([]byte(xml.Header))
	if err != nil {
		return err
	}

	b, err := xml.MarshalIndent(app, "", "\t")
	if err != nil {
		return err
	}

	_, err = replacer.WriteString(w, string(b))
	if err != nil {
		return err
	}

	return w.Close()
}
