package sparkle

import "encoding/xml"

type RSS struct {
	Channels []Channel `xml:"channel"`
}

func (r *RSS) MarshalXML(enc *xml.Encoder, _ xml.StartElement) error {
	err := enc.EncodeToken(xml.StartElement{
		Name: xml.Name{Local: "rss"},
		Attr: []xml.Attr{
			{Name: xml.Name{Local: "version"}, Value: "2.0"},
			{Name: xml.Name{Local: "xmlns:sparkle"}, Value: "http://www.andymatuschak.org/xml-namespaces/sparkle"},
			{Name: xml.Name{Local: "xmlns:dc"}, Value: "http://purl.org/dc/elements/1.1/"},
		},
	})
	if err != nil {
		return err
	}

	err = enc.EncodeElement(r.Channels, xml.StartElement{Name: xml.Name{Local: "channel"}})
	if err != nil {
		return err
	}

	return enc.EncodeToken(xml.EndElement{Name: xml.Name{Local: "rss"}})
}

func (r *RSS) UnmarshalXML(dec *xml.Decoder, _ xml.StartElement) error {
	var data []unmarshalRSS
	if err := dec.Decode(&data); err != nil {
		return err
	}

	r.Channels = make([]Channel, 0, len(data))
	for _, item := range data {
		channel := Channel{
			Title:       item.Title,
			Link:        item.Link,
			Description: item.Description,
			Language:    item.Language,
			Items:       make([]Item, 0, len(item.Items)),
		}

		for _, item := range item.Items {
			channel.Items = append(channel.Items, Item{
				Title:                             item.Title,
				PubDate:                           item.PubDate,
				Description:                       item.Description,
				Version:                           item.Version,
				ReleaseNotesLink:                  item.ReleaseNotesLink,
				CriticalUpdate:                    (*CriticalUpdate)(item.CriticalUpdate),
				Tags:                              (*Tags)(item.Tags),
				MinimumAutoupdateVersion:          item.MinimumAutoupdateVersion,
				IgnoreSkippedUpgradesBelowVersion: item.IgnoreSkippedUpgradesBelowVersion,
				Enclosure:                         Enclosure(item.Enclosure),
			})
		}

		r.Channels = append(r.Channels, channel)
	}

	return dec.Skip()
}

type Channel struct {
	Title       string `xml:"title"`
	Link        string `xml:"link,omitempty"`
	Description string `xml:"description,omitempty"`
	Language    string `xml:"language,omitempty"`
	Items       []Item `xml:"item"`
}

type Item struct {
	Title                             string          `xml:"title"`
	PubDate                           string          `xml:"pubDate"`
	Description                       *CdataString    `xml:"description,omitempty"`
	Version                           string          `xml:"sparkle:version,omitempty"`
	ReleaseNotesLink                  string          `xml:"sparkle:releaseNotesLink,omitempty"`
	CriticalUpdate                    *CriticalUpdate `xml:"sparkle:criticalUpdate,omitempty"`
	Tags                              *Tags           `xml:"sparkle:tags,omitempty"`
	MinimumAutoupdateVersion          string          `xml:"sparkle:minimumAutoupdateVersion,omitempty"`
	IgnoreSkippedUpgradesBelowVersion string          `xml:"sparkle:ignoreSkippedUpgradesBelowVersion,omitempty"`
	Enclosure                         Enclosure       `xml:"enclosure,omitempty"`
}

// CdataString for XML CDATA
// See issue: https://github.com/golang/go/issues/16198
type CdataString struct {
	Value string `xml:",cdata"`
}

type CriticalUpdate struct {
	Version string `xml:"sparkle:version,attr,omitempty"`
}

type Tags struct {
	CriticalUpdate Bool `xml:"sparkle:criticalUpdate,omitempty"`
}

type Bool bool

func (xb *Bool) MarshalText() ([]byte, error) {
	return nil, nil
}

func (xb *Bool) UnmarshalText([]byte) error {
	*xb = true
	return nil
}

type Enclosure struct {
	URL                  string `xml:"url,attr"`
	OS                   string `xml:"sparkle:os,attr"`
	Version              string `xml:"sparkle:version,attr"`
	DSASignature         string `xml:"sparkle:dsaSignature,attr,omitempty"`
	EDSignature          string `xml:"sparkle:edSignature,attr,omitempty"`
	InstallerArguments   string `xml:"sparkle:installerArguments,attr,omitempty"`
	MinimumSystemVersion string `xml:"sparkle:minimumSystemVersion,attr,omitempty"`
	Length               int    `xml:"length,attr,omitempty"`
	Type                 string `xml:"type,attr"`
}

type unmarshalRSS struct {
	XMLName     xml.Name `xml:"channel"`
	Title       string   `xml:"title"`
	Link        string   `xml:"link,omitempty"`
	Description string   `xml:"description,omitempty"`
	Language    string   `xml:"language,omitempty"`
	Items       []struct {
		Title            string       `xml:"title"`
		PubDate          string       `xml:"pubDate"`
		Description      *CdataString `xml:"description,omitempty"`
		Version          string       `xml:"version,omitempty"`
		ReleaseNotesLink string       `xml:"releaseNotesLink,omitempty"`
		CriticalUpdate   *struct {
			Version string `xml:"version,attr,omitempty"`
		} `xml:"criticalUpdate,omitempty"`
		Tags *struct {
			CriticalUpdate Bool `xml:"criticalUpdate,omitempty"`
		} `xml:"tags,omitempty"`
		MinimumAutoupdateVersion          string `xml:"minimumAutoupdateVersion,omitempty"`
		IgnoreSkippedUpgradesBelowVersion string `xml:"ignoreSkippedUpgradesBelowVersion,omitempty"`
		Enclosure                         struct {
			URL                  string `xml:"url,attr"`
			OS                   string `xml:"os,attr"`
			Version              string `xml:"version,attr"`
			DSASignature         string `xml:"dsaSignature,attr,omitempty"`
			EDSignature          string `xml:"edSignature,attr,omitempty"`
			InstallerArguments   string `xml:"installerArguments,attr,omitempty"`
			MinimumSystemVersion string `xml:"minimumSystemVersion,attr,omitempty"`
			Length               int    `xml:"length,attr,omitempty"`
			Type                 string `xml:"type,attr"`
		} `xml:"enclosure,omitempty"`
	} `xml:"item"`
}
