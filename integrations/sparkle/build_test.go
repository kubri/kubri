package sparkle_test

import (
	"context"
	"crypto/sha1"
	"encoding/base64"
	"encoding/xml"
	"io"
	"testing"
	"time"

	"github.com/abemedia/appcast/integrations/sparkle"
	"github.com/abemedia/appcast/internal/testsource"
	"github.com/abemedia/appcast/pkg/crypto/dsa"
	"github.com/abemedia/appcast/pkg/crypto/ed25519"
	"github.com/abemedia/appcast/source"
	target "github.com/abemedia/appcast/target/file"
	"github.com/google/go-cmp/cmp"
)

func TestBuild(t *testing.T) {
	ctx := context.Background()
	data := []byte("test")
	ts := time.Now().UTC()
	src := testsource.New([]*source.Release{
		{
			Version: "v1.0.0",
			Date:    ts,
		},
		{
			Version: "v1.1.0",
			Date:    ts,
			Description: `## New Features
- Something
- Something else`,
		},
	})
	src.UploadAsset(ctx, "v1.0.0", "test.dmg", data)
	src.UploadAsset(ctx, "v1.0.0", "test_32-bit.exe", data)
	src.UploadAsset(ctx, "v1.0.0", "test_64-bit.msi", data)
	src.UploadAsset(ctx, "v1.1.0", "test.dmg", data)
	src.UploadAsset(ctx, "v1.1.0", "test_32-bit.exe", data)
	src.UploadAsset(ctx, "v1.1.0", "test_64-bit.msi", data)

	tgt, err := target.New(target.Config{Path: t.TempDir(), URL: "https://example.com"})
	if err != nil {
		t.Fatal(err)
	}

	w, err := tgt.NewWriter(ctx, "appcast.xml")
	if err != nil {
		t.Fatal(err)
	}
	w.Write([]byte(`<?xml version="1.0" encoding="UTF-8"?>
<rss version="2.0" xmlns:sparkle="http://www.andymatuschak.org/xml-namespaces/sparkle" xmlns:dc="http://purl.org/dc/elements/1.1/">
	<channel>
		<item>
			<title>v1.0.0</title>
			<pubDate>Mon, 02 Jan 2006 15:04:05 +0000</pubDate>
			<sparkle:version>1.0.0</sparkle:version>
			<enclosure url="https://example.com/v1.0.0/test.dmg" sparkle:os="macos" sparkle:version="1.0.0" sparkle:edSignature="test" sparkle:minimumSystemVersion="10.13.0" length="4" type="application/x-apple-diskimage" />
		</item>
		<item>
			<title>v1.0.0</title>
			<pubDate>Mon, 02 Jan 2006 15:04:05 +0000</pubDate>
			<sparkle:version>1.0.0</sparkle:version>
			<sparkle:tags>
				<sparkle:criticalUpdate />
			</sparkle:tags>
			<enclosure url="https://example.com/v1.0.0/test_32-bit.exe" sparkle:os="windows-x86" sparkle:version="1.0.0" sparkle:dsaSignature="test" sparkle:installerArguments="/passive" length="4" type="application/vnd.microsoft.portable-executable" />
		</item>
		<item>
			<title>v1.0.0</title>
			<pubDate>Mon, 02 Jan 2006 15:04:05 +0000</pubDate>
			<sparkle:version>1.0.0</sparkle:version>
			<sparkle:tags>
				<sparkle:criticalUpdate />
			</sparkle:tags>
			<enclosure url="https://example.com/v1.0.0/test_64-bit.msi" sparkle:os="windows-x64" sparkle:version="1.0.0" sparkle:dsaSignature="test" sparkle:installerArguments="/passive" length="4" type="application/x-msi" />
		</item>
	</channel>
</rss>`))
	w.Close()

	c := &sparkle.Config{
		Title:       "Test",
		Description: "Test",
		URL:         "https://example.com/appcast.xml",
		Source:      src,
		Target:      tgt,
		FileName:    "appcast.xml",
		Settings: []sparkle.Rule{
			{
				OS: sparkle.Windows,
				Settings: &sparkle.Settings{
					InstallerArguments: "/passive",
				},
			},
			{
				OS: sparkle.MacOS,
				Settings: &sparkle.Settings{
					MinimumSystemVersion: "10.13.0",
				},
			},
			{
				Version: "v1.0",
				Settings: &sparkle.Settings{
					CriticalUpdate: true,
				},
			},
			{
				Version: ">= v1.1",
				Settings: &sparkle.Settings{
					CriticalUpdateBelowVersion:        "1.0.0",
					MinimumAutoupdateVersion:          "1.0.0",
					IgnoreSkippedUpgradesBelowVersion: "1.0.0",
				},
			},
		},
	}

	pubDate := ts.Format(time.RFC1123)

	want := `<?xml version="1.0" encoding="UTF-8"?>
<rss version="2.0" xmlns:sparkle="http://www.andymatuschak.org/xml-namespaces/sparkle" xmlns:dc="http://purl.org/dc/elements/1.1/">
	<channel>
		<title>Test</title>
		<link>https://example.com/appcast.xml</link>
		<description>Test</description>
		<item>
			<title>v1.1.0</title>
			<pubDate>` + pubDate + `</pubDate>
			<description><![CDATA[
				<h2>New Features</h2>
				<ul>
					<li>Something</li>
					<li>Something else</li>
				</ul>
			]]></description>
			<sparkle:version>1.1.0</sparkle:version>
			<sparkle:criticalUpdate sparkle:version="1.0.0" />
			<sparkle:minimumAutoupdateVersion>1.0.0</sparkle:minimumAutoupdateVersion>
			<sparkle:ignoreSkippedUpgradesBelowVersion>1.0.0</sparkle:ignoreSkippedUpgradesBelowVersion>
			<enclosure url="https://example.com/v1.1.0/test.dmg" sparkle:os="macos" sparkle:version="1.1.0" sparkle:minimumSystemVersion="10.13.0" length="4" type="application/x-apple-diskimage" />
		</item>
		<item>
			<title>v1.1.0</title>
			<pubDate>` + pubDate + `</pubDate>
			<description><![CDATA[
				<h2>New Features</h2>
				<ul>
					<li>Something</li>
					<li>Something else</li>
				</ul>
			]]></description>
			<sparkle:version>1.1.0</sparkle:version>
			<sparkle:criticalUpdate sparkle:version="1.0.0" />
			<sparkle:minimumAutoupdateVersion>1.0.0</sparkle:minimumAutoupdateVersion>
			<sparkle:ignoreSkippedUpgradesBelowVersion>1.0.0</sparkle:ignoreSkippedUpgradesBelowVersion>
			<enclosure url="https://example.com/v1.1.0/test_32-bit.exe" sparkle:os="windows-x86" sparkle:version="1.1.0" sparkle:installerArguments="/passive" length="4" type="application/vnd.microsoft.portable-executable" />
		</item>
		<item>
			<title>v1.1.0</title>
			<pubDate>` + pubDate + `</pubDate>
			<description><![CDATA[
				<h2>New Features</h2>
				<ul>
					<li>Something</li>
					<li>Something else</li>
				</ul>
			]]></description>
			<sparkle:version>1.1.0</sparkle:version>
			<sparkle:criticalUpdate sparkle:version="1.0.0" />
			<sparkle:minimumAutoupdateVersion>1.0.0</sparkle:minimumAutoupdateVersion>
			<sparkle:ignoreSkippedUpgradesBelowVersion>1.0.0</sparkle:ignoreSkippedUpgradesBelowVersion>
			<enclosure url="https://example.com/v1.1.0/test_64-bit.msi" sparkle:os="windows-x64" sparkle:version="1.1.0" sparkle:installerArguments="/passive" length="4" type="application/x-msi" />
		</item>
		<item>
			<title>v1.0.0</title>
			<pubDate>Mon, 02 Jan 2006 15:04:05 +0000</pubDate>
			<sparkle:version>1.0.0</sparkle:version>
			<enclosure url="https://example.com/v1.0.0/test.dmg" sparkle:os="macos" sparkle:version="1.0.0" sparkle:edSignature="test" sparkle:minimumSystemVersion="10.13.0" length="4" type="application/x-apple-diskimage" />
		</item>
		<item>
			<title>v1.0.0</title>
			<pubDate>Mon, 02 Jan 2006 15:04:05 +0000</pubDate>
			<sparkle:version>1.0.0</sparkle:version>
			<sparkle:tags>
				<sparkle:criticalUpdate />
			</sparkle:tags>
			<enclosure url="https://example.com/v1.0.0/test_32-bit.exe" sparkle:os="windows-x86" sparkle:version="1.0.0" sparkle:dsaSignature="test" sparkle:installerArguments="/passive" length="4" type="application/vnd.microsoft.portable-executable" />
		</item>
		<item>
			<title>v1.0.0</title>
			<pubDate>Mon, 02 Jan 2006 15:04:05 +0000</pubDate>
			<sparkle:version>1.0.0</sparkle:version>
			<sparkle:tags>
				<sparkle:criticalUpdate />
			</sparkle:tags>
			<enclosure url="https://example.com/v1.0.0/test_64-bit.msi" sparkle:os="windows-x64" sparkle:version="1.0.0" sparkle:dsaSignature="test" sparkle:installerArguments="/passive" length="4" type="application/x-msi" />
		</item>
	</channel>
</rss>`

	testBuild(t, c, want)

	// Should be no-op as nothing changed so timestamp should still be valid.
	time.Sleep(time.Second)
	testBuild(t, c, want)
}

func testBuild(t *testing.T, c *sparkle.Config, want string) {
	t.Helper()

	ctx := context.Background()
	if err := sparkle.Build(ctx, c); err != nil {
		t.Fatal(err)
	}

	r, err := c.Target.NewReader(ctx, "appcast.xml")
	if err != nil {
		t.Fatal(err)
	}
	defer r.Close()

	got, err := io.ReadAll(r)
	if err != nil {
		t.Fatal(err)
	}

	if diff := cmp.Diff(want, string(got)); diff != "" {
		t.Error(diff)
	}
}

func TestBuildSign(t *testing.T) {
	ctx := context.Background()
	data := []byte("test")
	ts := time.Now().UTC()
	src := testsource.New([]*source.Release{{Version: "v1.0.0", Date: ts}})
	src.UploadAsset(ctx, "v1.0.0", "test.dmg", data)
	src.UploadAsset(ctx, "v1.0.0", "test.msi", data)

	tgt, err := target.New(target.Config{Path: t.TempDir(), URL: "https://example.com"})
	if err != nil {
		t.Fatal(err)
	}

	dsaKey, err := dsa.NewPrivateKey()
	if err != nil {
		t.Fatal(err)
	}

	edKey, err := ed25519.NewPrivateKey()
	if err != nil {
		t.Fatal(err)
	}

	c := &sparkle.Config{
		Title:       "Test",
		Description: "Test",
		URL:         "https://example.com/appcast.xml",
		Source:      src,
		Target:      tgt,
		FileName:    "appcast.xml",
		Settings:    []sparkle.Rule{},
		DSAKey:      dsaKey,
		Ed25519Key:  edKey,
	}

	pubDate := ts.Format(time.RFC1123)

	want := sparkle.RSS{
		Channels: []*sparkle.Channel{{
			Title:       "Test",
			Link:        "https://example.com/appcast.xml",
			Description: "Test",
			Items: []*sparkle.Item{
				{
					Title:   "v1.0.0",
					PubDate: pubDate,
					Version: "1.0.0",
					Enclosure: &sparkle.Enclosure{
						URL:         "https://example.com/v1.0.0/test.dmg",
						OS:          "macos",
						Version:     "1.0.0",
						EDSignature: base64.RawStdEncoding.EncodeToString(ed25519.Sign(edKey, data)),
						Length:      4,
						Type:        "application/x-apple-diskimage",
					},
				},
				{
					Title:   "v1.0.0",
					PubDate: pubDate,
					Version: "1.0.0",
					Enclosure: &sparkle.Enclosure{
						URL:     "https://example.com/v1.0.0/test.msi",
						OS:      "windows",
						Version: "1.0.0",
						DSASignature: func() string {
							sum := sha1.Sum(data)
							sum = sha1.Sum(sum[:])
							sig, _ := dsa.Sign(dsaKey, sum[:])
							return base64.RawStdEncoding.EncodeToString(sig)
						}(),
						Length: 4,
						Type:   "application/x-msi",
					},
				},
			},
		}},
	}

	if err = sparkle.Build(ctx, c); err != nil {
		t.Fatal(err)
	}

	r, err := tgt.NewReader(ctx, "appcast.xml")
	if err != nil {
		t.Fatal(err)
	}
	defer r.Close()

	var got sparkle.RSS
	if err = xml.NewDecoder(r).Decode(&got); err != nil {
		t.Fatal(err)
	}

	compareDSA := cmp.FilterPath(func(p cmp.Path) bool {
		return p.String() == "Channels.Items.Enclosure.DSASignature"
	}, cmp.Comparer(func(a, b string) bool {
		if a == "" || b == "" {
			return a == b
		}
		x, _ := base64.RawStdEncoding.DecodeString(a)
		y, _ := base64.RawStdEncoding.DecodeString(b)
		pub := dsa.Public(dsaKey)
		sum := sha1.Sum(data)
		sum = sha1.Sum(sum[:])
		return dsa.Verify(pub, sum[:], x) && dsa.Verify(pub, sum[:], y)
	}))

	compareED := cmp.FilterPath(func(p cmp.Path) bool {
		return p.String() == "Channels.Items.Enclosure.EDSignature"
	}, cmp.Comparer(func(a, b string) bool {
		if a == "" || b == "" {
			return a == b
		}
		x, _ := base64.RawStdEncoding.DecodeString(a)
		y, _ := base64.RawStdEncoding.DecodeString(b)
		pub := ed25519.Public(edKey)
		return ed25519.Verify(pub, data, x) && ed25519.Verify(pub, data, y)
	}))

	if diff := cmp.Diff(want, got, compareDSA, compareED); diff != "" {
		t.Error(diff)
	}
}

func TestBuildUpload(t *testing.T) {
	ctx := context.Background()
	data := []byte("test")
	ts := time.Now().UTC()
	src := testsource.New([]*source.Release{{Version: "v1.0.0", Date: ts}})
	src.UploadAsset(ctx, "v1.0.0", "test.dmg", data)
	src.UploadAsset(ctx, "v1.0.0", "test.msi", data)

	for _, upload := range []bool{true, false} {
		tgt, err := target.New(target.Config{Path: t.TempDir()})
		if err != nil {
			t.Fatal(err)
		}

		c := &sparkle.Config{
			Title:          "Test",
			Description:    "Test",
			URL:            "https://example.com/appcast.xml",
			Source:         src,
			Target:         tgt,
			FileName:       "appcast.xml",
			Settings:       []sparkle.Rule{},
			UploadPackages: upload,
		}

		if err := sparkle.Build(ctx, c); err != nil {
			t.Fatal(err)
		}

		r, err := src.GetRelease(ctx, "v1.0.0")
		if err != nil {
			t.Fatal(err)
		}

		for _, asset := range r.Assets {
			rd, err := tgt.NewReader(ctx, r.Version+"/"+asset.Name)
			if err == nil {
				rd.Close()
			}
			if upload && err != nil {
				t.Fatalf("should upload assets: %v", err)
			}
			if !upload && err == nil {
				t.Fatal("should not upload assets")
			}
		}
	}
}
