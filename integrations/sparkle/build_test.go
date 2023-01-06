package sparkle_test

import (
	"context"
	"io"
	"testing"
	"time"

	"github.com/abemedia/appcast/integrations/sparkle"
	"github.com/abemedia/appcast/source"
	"github.com/abemedia/appcast/source/blob/memory"
	target "github.com/abemedia/appcast/target/blob/memory"
	"github.com/google/go-cmp/cmp"
)

func TestBuild(t *testing.T) {
	data := []byte("test")
	ctx := context.Background()

	src, _ := memory.New(source.Config{})
	src.UploadAsset(ctx, "v1.0.0", "test.dmg", data)
	src.UploadAsset(ctx, "v1.0.0", "test_64-bit.msi", data)
	src.UploadAsset(ctx, "v1.0.0", "test_32-bit.exe", data)
	src.UploadAsset(ctx, "v1.1.0", "test.dmg", data)
	src.UploadAsset(ctx, "v1.1.0", "test_64-bit.msi", data)
	src.UploadAsset(ctx, "v1.1.0", "test_32-bit.exe", data)

	tgt, err := target.New(source.Config{})
	if err != nil {
		t.Fatal(err)
	}

	w, err := tgt.NewWriter(ctx, "sparkle.xml")
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
			<enclosure url="mem://v1.0.0/test.dmg" sparkle:os="macos" sparkle:version="1.0.0" sparkle:edSignature="test" sparkle:minimumSystemVersion="10.13.0" length="4" type="application/x-apple-diskimage" />
		</item>
		<item>
			<title>v1.0.0</title>
			<pubDate>Mon, 02 Jan 2006 15:04:05 +0000</pubDate>
			<sparkle:version>1.0.0</sparkle:version>
			<sparkle:tags>
				<sparkle:criticalUpdate />
			</sparkle:tags>
			<enclosure url="mem://v1.0.0/test_32-bit.exe" sparkle:os="windows-x86" sparkle:version="1.0.0" sparkle:dsaSignature="test" sparkle:installerArguments="/passive" length="4" type="application/vnd.microsoft.portable-executable" />
		</item>
		<item>
			<title>v1.0.0</title>
			<pubDate>Mon, 02 Jan 2006 15:04:05 +0000</pubDate>
			<sparkle:version>1.0.0</sparkle:version>
			<sparkle:tags>
				<sparkle:criticalUpdate />
			</sparkle:tags>
			<enclosure url="mem://v1.0.0/test_64-bit.msi" sparkle:os="windows-x64" sparkle:version="1.0.0" sparkle:dsaSignature="test" sparkle:installerArguments="/passive" length="4" type="application/x-msi" />
		</item>
	</channel>
</rss>`))
	w.Close()

	c := &sparkle.Config{
		Title:       "Test",
		Description: "Test",
		URL:         "https://example.com/sparkle.xml",
		Source:      src,
		Target:      tgt,
		FileName:    "sparkle.xml",
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

	pubDate := time.Now().UTC().Format(time.RFC1123)

	want := `<?xml version="1.0" encoding="UTF-8"?>
<rss version="2.0" xmlns:sparkle="http://www.andymatuschak.org/xml-namespaces/sparkle" xmlns:dc="http://purl.org/dc/elements/1.1/">
	<channel>
		<title>Test</title>
		<link>https://example.com/sparkle.xml</link>
		<description>Test</description>
		<item>
			<title>v1.1.0</title>
			<pubDate>` + pubDate + `</pubDate>
			<sparkle:version>1.1.0</sparkle:version>
			<sparkle:criticalUpdate sparkle:version="1.0.0" />
			<sparkle:minimumAutoupdateVersion>1.0.0</sparkle:minimumAutoupdateVersion>
			<sparkle:ignoreSkippedUpgradesBelowVersion>1.0.0</sparkle:ignoreSkippedUpgradesBelowVersion>
			<enclosure url="mem://v1.1.0/test.dmg" sparkle:os="macos" sparkle:version="1.1.0" sparkle:minimumSystemVersion="10.13.0" length="4" type="application/x-apple-diskimage" />
		</item>
		<item>
			<title>v1.1.0</title>
			<pubDate>` + pubDate + `</pubDate>
			<sparkle:version>1.1.0</sparkle:version>
			<sparkle:criticalUpdate sparkle:version="1.0.0" />
			<sparkle:minimumAutoupdateVersion>1.0.0</sparkle:minimumAutoupdateVersion>
			<sparkle:ignoreSkippedUpgradesBelowVersion>1.0.0</sparkle:ignoreSkippedUpgradesBelowVersion>
			<enclosure url="mem://v1.1.0/test_32-bit.exe" sparkle:os="windows-x86" sparkle:version="1.1.0" sparkle:installerArguments="/passive" length="4" type="application/vnd.microsoft.portable-executable" />
		</item>
		<item>
			<title>v1.1.0</title>
			<pubDate>` + pubDate + `</pubDate>
			<sparkle:version>1.1.0</sparkle:version>
			<sparkle:criticalUpdate sparkle:version="1.0.0" />
			<sparkle:minimumAutoupdateVersion>1.0.0</sparkle:minimumAutoupdateVersion>
			<sparkle:ignoreSkippedUpgradesBelowVersion>1.0.0</sparkle:ignoreSkippedUpgradesBelowVersion>
			<enclosure url="mem://v1.1.0/test_64-bit.msi" sparkle:os="windows-x64" sparkle:version="1.1.0" sparkle:installerArguments="/passive" length="4" type="application/x-msi" />
		</item>
		<item>
			<title>v1.0.0</title>
			<pubDate>Mon, 02 Jan 2006 15:04:05 +0000</pubDate>
			<sparkle:version>1.0.0</sparkle:version>
			<enclosure url="mem://v1.0.0/test.dmg" sparkle:os="macos" sparkle:version="1.0.0" sparkle:edSignature="test" sparkle:minimumSystemVersion="10.13.0" length="4" type="application/x-apple-diskimage" />
		</item>
		<item>
			<title>v1.0.0</title>
			<pubDate>Mon, 02 Jan 2006 15:04:05 +0000</pubDate>
			<sparkle:version>1.0.0</sparkle:version>
			<sparkle:tags>
				<sparkle:criticalUpdate />
			</sparkle:tags>
			<enclosure url="mem://v1.0.0/test_32-bit.exe" sparkle:os="windows-x86" sparkle:version="1.0.0" sparkle:dsaSignature="test" sparkle:installerArguments="/passive" length="4" type="application/vnd.microsoft.portable-executable" />
		</item>
		<item>
			<title>v1.0.0</title>
			<pubDate>Mon, 02 Jan 2006 15:04:05 +0000</pubDate>
			<sparkle:version>1.0.0</sparkle:version>
			<sparkle:tags>
				<sparkle:criticalUpdate />
			</sparkle:tags>
			<enclosure url="mem://v1.0.0/test_64-bit.msi" sparkle:os="windows-x64" sparkle:version="1.0.0" sparkle:dsaSignature="test" sparkle:installerArguments="/passive" length="4" type="application/x-msi" />
		</item>
	</channel>
</rss>`

	err = sparkle.Build(ctx, c)
	if err != nil {
		t.Fatal(err)
	}

	r, err := tgt.NewReader(ctx, "sparkle.xml")
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
