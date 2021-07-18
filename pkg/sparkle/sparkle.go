package sparkle

import "encoding/xml"

type Feed struct {
	XMLName      xml.Name `xml:"rss"`
	Version      string   `xml:"version,attr"`
	XMLNSSparkle string   `xml:"xmlns:sparkle,attr"`
	XMLNSDC      string   `xml:"xmlns:dc,attr"`
	Channels     []Channel
}

type Channel struct {
	XMLName     xml.Name `xml:"channel"`
	Title       string   `xml:"title"`
	Link        string   `xml:"link,omitempty"`
	Description string   `xml:"description,omitempty"`
	Language    string   `xml:"language,omitempty"`
	Items       []Item
}

type Item struct {
	XMLName                           xml.Name     `xml:"item"`
	Title                             string       `xml:"title"`
	Version                           string       `xml:"sparkle:version,omitempty"`
	ReleaseNotesLink                  string       `xml:"sparkle:releaseNotesLink,omitempty"`
	MinimumAutoupdateVersion          string       `xml:"sparkle:minimumAutoupdateVersion,omitempty"`
	IgnoreSkippedUpgradesBelowVersion string       `xml:"sparkle:ignoreSkippedUpgradesBelowVersion,omitempty"`
	Description                       *CdataString `xml:"description,omitempty"`
	PubDate                           string       `xml:"pubDate"`
	Enclosure                         Enclosure    `xml:"enclosure,omitempty"`
}

// CdataString for XML CDATA
// See issue: https://github.com/golang/go/issues/16198
type CdataString struct {
	Value string `xml:",cdata"`
}

type Enclosure struct {
	XMLName              xml.Name `xml:"enclosure"`
	URL                  string   `xml:"url,attr"`
	OS                   string   `xml:"sparkle:os,attr"`
	Version              string   `xml:"sparkle:version,attr"`
	DSASignature         string   `xml:"sparkle:dsaSignature,attr,omitempty"`
	EDSignature          string   `xml:"sparkle:edSignature,attr,omitempty"`
	InstallerArguments   string   `xml:"sparkle:installerArguments,attr,omitempty"`
	MinimumSystemVersion string   `xml:"sparkle:minimumSystemVersion,attr,omitempty"`
	Length               int      `xml:"length,attr,omitempty"`
	Type                 string   `xml:"type,attr"`
}
