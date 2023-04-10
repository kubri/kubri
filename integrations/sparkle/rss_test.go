package sparkle_test

import (
	"encoding/xml"
	"testing"

	"github.com/abemedia/appcast/integrations/sparkle"
	"github.com/google/go-cmp/cmp"
)

func TestRSSMarshalUnmarshal(t *testing.T) {
	in := &sparkle.RSS{
		Channels: []*sparkle.Channel{{
			Title:       "Title",
			Description: "Description",
			Link:        "https://example.com/sparkle.xml",
			Language:    "en-gb",
			Items: []*sparkle.Item{
				{
					Title:       "v1.0.0",
					Description: &sparkle.CdataString{"Test"},
					PubDate:     "Mon, 02 Jan 2006 15:04:05 +0000",
					Version:     "1.0.0",
					Tags:        &sparkle.Tags{CriticalUpdate: true},
					Enclosure: &sparkle.Enclosure{
						URL:         "https://example.com/test_v1.0.0.dmg",
						OS:          "macos",
						Version:     "1.0.0",
						EDSignature: "test",
						Length:      100,
						Type:        "application/x-apple-diskimage",
					},
				},
				{
					Title:          "v1.1.0",
					Description:    &sparkle.CdataString{"\n\t\t\t\t<h2>Test</h2>\n\t\t\t"},
					PubDate:        "Mon, 02 Jan 2007 15:04:05 +0000",
					Version:        "1.1.0",
					CriticalUpdate: &sparkle.CriticalUpdate{Version: "1.0.0"},
					Enclosure: &sparkle.Enclosure{
						URL:         "https://example.com/test_v1.1.0.dmg",
						OS:          "macos",
						Version:     "1.1.0",
						EDSignature: "test",
						Length:      100,
						Type:        "application/x-apple-diskimage",
					},
				},
			},
		}},
	}

	want := `<rss version="2.0" xmlns:sparkle="http://www.andymatuschak.org/xml-namespaces/sparkle" xmlns:dc="http://purl.org/dc/elements/1.1/">
	<channel>
		<title>Title</title>
		<link>https://example.com/sparkle.xml</link>
		<description>Description</description>
		<language>en-gb</language>
		<item>
			<title>v1.0.0</title>
			<pubDate>Mon, 02 Jan 2006 15:04:05 +0000</pubDate>
			<description><![CDATA[Test]]></description>
			<sparkle:version>1.0.0</sparkle:version>
			<sparkle:tags>
				<sparkle:criticalUpdate></sparkle:criticalUpdate>
			</sparkle:tags>
			<enclosure url="https://example.com/test_v1.0.0.dmg" sparkle:os="macos" sparkle:version="1.0.0" sparkle:edSignature="test" length="100" type="application/x-apple-diskimage"></enclosure>
		</item>
		<item>
			<title>v1.1.0</title>
			<pubDate>Mon, 02 Jan 2007 15:04:05 +0000</pubDate>
			<description><![CDATA[
				<h2>Test</h2>
			]]></description>
			<sparkle:version>1.1.0</sparkle:version>
			<sparkle:criticalUpdate sparkle:version="1.0.0"></sparkle:criticalUpdate>
			<enclosure url="https://example.com/test_v1.1.0.dmg" sparkle:os="macos" sparkle:version="1.1.0" sparkle:edSignature="test" length="100" type="application/x-apple-diskimage"></enclosure>
		</item>
	</channel>
</rss>`

	b, err := xml.MarshalIndent(in, "", "\t")
	if err != nil {
		t.Fatal(err)
	}

	if diff := cmp.Diff(want, string(b)); diff != "" {
		t.Fatal(diff)
	}

	got := &sparkle.RSS{}
	if err = xml.Unmarshal(b, got); err != nil {
		t.Fatal(err)
	}

	if diff := cmp.Diff(in, got); diff != "" {
		t.Fatal(diff)
	}
}
