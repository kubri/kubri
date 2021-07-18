package sparkle_test

import (
	"encoding/xml"
	"testing"

	"github.com/abemedia/appcast/pkg/sparkle"
	"github.com/google/go-cmp/cmp"
)

func TestSparkleXMLMarshal(t *testing.T) {
	in := sparkle.Feed{
		Version:      "2.0",
		XMLNSSparkle: "http://www.andymatuschak.org/xml-namespaces/sparkle",
		XMLNSDC:      "http://purl.org/dc/elements/1.1/",
		Channels: []sparkle.Channel{
			{
				Title:       "Test",
				Link:        "https://www.example.com",
				Description: "Test",
				Items: []sparkle.Item{
					{
						Title:                             "Version 1.1.0",
						PubDate:                           "Mon, 20 Oct 2015 19:20:11 +0000",
						Description:                       &sparkle.CdataString{Value: "This is v1"},
						Version:                           "1.1.0",
						ReleaseNotesLink:                  "https://example.com/release-notes",
						CriticalUpdate:                    &sparkle.CriticalUpdate{Version: "1.0.0"},
						MinimumAutoupdateVersion:          "1.0.0",
						IgnoreSkippedUpgradesBelowVersion: "1.0.0",
						Enclosure: sparkle.Enclosure{
							URL:                "https://example.com/test.msi",
							OS:                 "windows",
							Version:            "1.1.0",
							DSASignature:       "test",
							InstallerArguments: "/passive",
							Length:             1000,
							Type:               "application/octet-stream",
						},
					},
					{
						Title:            "Version 1.0.0",
						PubDate:          "Mon, 05 Oct 2015 19:20:11 +0000",
						Description:      &sparkle.CdataString{Value: "This is v1"},
						Version:          "1.0.0",
						ReleaseNotesLink: "https://example.com/release-notes",
						Tags:             &sparkle.Tags{CriticalUpdate: true},
						Enclosure: sparkle.Enclosure{
							URL:                  "https://example.com/test.dmg",
							OS:                   "macos",
							Version:              "1.0.0",
							EDSignature:          "test",
							MinimumSystemVersion: "10.13.0",
							Length:               1000,
							Type:                 "application/octet-stream",
						},
					},
				},
			},
		},
	}

	want := `<rss version="2.0" xmlns:sparkle="http://www.andymatuschak.org/xml-namespaces/sparkle" xmlns:dc="http://purl.org/dc/elements/1.1/">
	<channel>
		<title>Test</title>
		<link>https://www.example.com</link>
		<description>Test</description>
		<item>
			<title>Version 1.1.0</title>
			<pubDate>Mon, 20 Oct 2015 19:20:11 +0000</pubDate>
			<description><![CDATA[This is v1]]></description>
			<sparkle:version>1.1.0</sparkle:version>
			<sparkle:releaseNotesLink>https://example.com/release-notes</sparkle:releaseNotesLink>
			<sparkle:criticalUpdate sparkle:version="1.0.0"></sparkle:criticalUpdate>
			<sparkle:minimumAutoupdateVersion>1.0.0</sparkle:minimumAutoupdateVersion>
			<sparkle:ignoreSkippedUpgradesBelowVersion>1.0.0</sparkle:ignoreSkippedUpgradesBelowVersion>
			<enclosure url="https://example.com/test.msi" sparkle:os="windows" sparkle:version="1.1.0" sparkle:dsaSignature="test" sparkle:installerArguments="/passive" length="1000" type="application/octet-stream"></enclosure>
		</item>
		<item>
			<title>Version 1.0.0</title>
			<pubDate>Mon, 05 Oct 2015 19:20:11 +0000</pubDate>
			<description><![CDATA[This is v1]]></description>
			<sparkle:version>1.0.0</sparkle:version>
			<sparkle:releaseNotesLink>https://example.com/release-notes</sparkle:releaseNotesLink>
			<sparkle:tags>
				<sparkle:criticalUpdate></sparkle:criticalUpdate>
			</sparkle:tags>
			<enclosure url="https://example.com/test.dmg" sparkle:os="macos" sparkle:version="1.0.0" sparkle:edSignature="test" sparkle:minimumSystemVersion="10.13.0" length="1000" type="application/octet-stream"></enclosure>
		</item>
	</channel>
</rss>`

	b, err := xml.MarshalIndent(in, "", "\t")
	if err != nil {
		t.Fatal(err)
	}

	if diff := cmp.Diff(want, string(b)); diff != "" {
		t.Error(diff)
	}
}
