package sparkle_test

import (
	"encoding/xml"
	"testing"

	"github.com/abemedia/appcast/integrations/sparkle"
	"github.com/google/go-cmp/cmp"
)

func TestRSSMarshalUnmarshal(t *testing.T) {
	in := &sparkle.RSS{
		Channels: []sparkle.Channel{{
			Title:       "Title",
			Description: "Description",
			Link:        "https://example.com/sparkle.xml",
			Language:    "en-gb",
			Items: []sparkle.Item{
				{
					Title:       "v1.0.0",
					Description: &sparkle.CdataString{"Test"},
					PubDate:     "Mon, 02 Jan 2006 15:04:05 +0000",
					Version:     "1.0.0",
					Tags:        &sparkle.Tags{CriticalUpdate: true},
					Enclosure: sparkle.Enclosure{
						URL:         "https://example.com/test_v1.0.0.dmg",
						OS:          "macos",
						Version:     "1.0.0",
						EDSignature: "test",
						Length:      100,
					},
				},
				{
					Title: "v1.1.0",
					Description: &sparkle.CdataString{`
						<h2>Test</h2>
					`},
					PubDate:        "Mon, 02 Jan 2007 15:04:05 +0000",
					Version:        "1.1.0",
					CriticalUpdate: &sparkle.CriticalUpdate{Version: "1.0.0"},
					Enclosure: sparkle.Enclosure{
						URL:         "https://example.com/test_v1.1.0.dmg",
						OS:          "macos",
						Version:     "1.1.0",
						EDSignature: "test",
						Length:      100,
					},
				},
			},
		}},
	}

	b, err := xml.Marshal(in)
	if err != nil {
		t.Fatal(err)
	}

	got := &sparkle.RSS{}
	if err = xml.Unmarshal(b, got); err != nil {
		t.Fatal(err)
	}

	if diff := cmp.Diff(in, got); diff != "" {
		t.Fatal(err)
	}
}
