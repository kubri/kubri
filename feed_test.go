package appcast_test

import (
	"testing"
	"time"

	"github.com/abemedia/appcast"
	"github.com/abemedia/appcast/pkg/sparkle"
	"github.com/abemedia/appcast/source"
	"github.com/abemedia/appcast/source/blob/memory"
	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
)

func TestFeed(t *testing.T) {
	data := []byte("test")

	s, _ := memory.New(source.Config{})
	s.UploadAsset("v1.0.0", "README.md", data)
	s.UploadAsset("v1.0.0", "test.dmg", data)
	s.UploadAsset("v1.0.0", "test_64-bit.msi", data)
	s.UploadAsset("v1.0.0", "test_32-bit.exe", data)
	s.UploadAsset("v1.0.0-alpha1", "test.dmg", data)
	s.UploadAsset("v1.0.0-alpha1", "test_64-bit.msi", data)

	c := &appcast.Config{
		Title:       "test",
		Description: "test",
		URL:         "https://example.com/updates.xml",
		Source:      s,
		Settings: []appcast.Rule{
			{
				OS: appcast.Windows,
				Settings: &appcast.Settings{
					InstallerArguments: "/passive",
				},
			},
			{
				OS: appcast.MacOS,
				Settings: &appcast.Settings{
					MinimumSystemVersion: "10.13.0",
				},
			},
		},
	}

	want := &sparkle.Feed{
		Version:      "2.0",
		XMLNSSparkle: "http://www.andymatuschak.org/xml-namespaces/sparkle",
		XMLNSDC:      "http://purl.org/dc/elements/1.1/",
		Channels: []sparkle.Channel{
			{
				Title:       "test",
				Description: "test",
				Link:        "https://example.com/updates.xml",
				Items: []sparkle.Item{
					{
						Title:   "v1.0.0",
						Version: "1.0.0",
						PubDate: time.Now().UTC().Format(time.RFC1123),
						Enclosure: sparkle.Enclosure{
							URL:                  "mem://v1.0.0/test.dmg",
							OS:                   "macos",
							Version:              "1.0.0",
							MinimumSystemVersion: "10.13.0",
							Length:               4,
							Type:                 "application/x-apple-diskimage",
						},
					},
					{
						Title:   "v1.0.0",
						Version: "1.0.0",
						PubDate: time.Now().UTC().Format(time.RFC1123),
						Enclosure: sparkle.Enclosure{
							URL:                "mem://v1.0.0/test_64-bit.msi",
							OS:                 "windows-x64",
							Version:            "1.0.0",
							InstallerArguments: "/passive",
							Length:             4,
							Type:               "application/x-msi",
						},
					},
					{
						Title:   "v1.0.0",
						Version: "1.0.0",
						PubDate: time.Now().UTC().Format(time.RFC1123),
						Enclosure: sparkle.Enclosure{
							URL:                "mem://v1.0.0/test_32-bit.exe",
							OS:                 "windows-x86",
							Version:            "1.0.0",
							InstallerArguments: "/passive",
							Length:             4,
							Type:               "application/vnd.microsoft.portable-executable",
						},
					},
				},
			},
		},
	}

	got, err := appcast.Feed(c)
	if err != nil {
		t.Fatal(err)
	}

	opt := cmpopts.SortSlices(func(a, b sparkle.Item) bool {
		return a.Version+a.Enclosure.OS > b.Version+b.Enclosure.OS
	})

	if diff := cmp.Diff(want, got, opt); diff != "" {
		t.Error(diff)
	}
}
