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
	HoursBetweenUpdateChecks  int
	UpdateBlocksActivation    bool
	ShowPrompt                bool
	AutomaticBackgroundTask   bool
	ForceUpdateFromAnyVersion bool

	Source         *source.Source
	Target         target.Target
	Version        string
	Prerelease     bool
	UploadPackages bool
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
				ok = true
				if err := build(ctx, c, r.Version, asset); err != nil {
					return err
				}
			}
		}
		if ok {
			break // Only build for latest release containing supported files.
		}
	}

	return nil
}

func build(ctx context.Context, c *Config, version string, asset *source.Asset) error {
	b, err := c.Source.DownloadAsset(ctx, version, asset.Name)
	if err != nil {
		return err
	}

	p, err := getPackage(b)
	if err != nil {
		return err
	}

	if c.UploadPackages {
		p.URI, err = upload(ctx, c.Target, asset.Name, b)
		if err != nil {
			return err
		}
	} else {
		p.URI = asset.URL
	}

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
		return err
	}

	return write(ctx, c, res)
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

func upload(ctx context.Context, t target.Target, path string, data []byte) (string, error) {
	w, err := t.NewWriter(ctx, path)
	if err != nil {
		return "", err
	}
	if _, err = w.Write(data); err != nil {
		return "", err
	}
	if err = w.Close(); err != nil {
		return "", err
	}
	return t.URL(ctx, path)
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

//nolint:gochecknoglobals
var replacer = strings.NewReplacer("></Package>", " />", "></Bundle>", " />",
	"></MainPackage>", " />", "></MainBundle>", " />", "></OnLaunch>", " />",
	"></AutomaticBackgroundTask>", " />", "></ForceUpdateFromAnyVersion>", " />")

func write(ctx context.Context, c *Config, app *AppInstaller) error {
	w, err := c.Target.NewWriter(ctx, path.Base(app.URI))
	if err != nil {
		return err
	}
	defer w.Close()
	if _, err = w.Write([]byte(xml.Header)); err != nil {
		return err
	}
	b, err := xml.MarshalIndent(app, "", "\t")
	if err != nil {
		return err
	}
	if _, err = replacer.WriteString(w, string(b)); err != nil {
		return err
	}
	return w.Close()
}
