package appinstaller

import "encoding/xml"

type AppInstaller struct {
	XMLName          xml.Name
	Version          string `xml:",attr"`
	URI              string `xml:"Uri,attr"`
	MainBundle       *Package
	MainPackage      *Package
	Dependencies     *Packages
	OptionalPackages *Packages
	RelatedPackages  *Packages
	UpdateSettings   *UpdateSettings
}

func (ai *AppInstaller) MarshalXML(enc *xml.Encoder, start xml.StartElement) error {
	start.Name.Local = "AppInstaller"
	start.Name.Space = ai.XMLName.Space
	return enc.EncodeElement(*ai, start)
}

type Packages struct {
	Bundle  []*Package
	Package []*Package
}

type Package struct {
	Name                  string `xml:",attr"`
	Publisher             string `xml:",attr"`
	Version               string `xml:",attr"`
	ProcessorArchitecture string `xml:",attr,omitempty"`
	URI                   string `xml:"Uri,attr"`
}

type UpdateSettings struct {
	OnLaunch                  *OnLaunch
	AutomaticBackgroundTask   Bool `xml:",omitempty"`
	ForceUpdateFromAnyVersion bool `xml:",omitempty"`
}

type OnLaunch struct {
	HoursBetweenUpdateChecks int  `xml:",attr,omitempty"`
	UpdateBlocksActivation   bool `xml:",attr,omitempty"`
	ShowPrompt               bool `xml:",attr,omitempty"`
}

type Bool bool

func (xb *Bool) MarshalText() ([]byte, error) {
	return nil, nil
}

func (xb *Bool) UnmarshalText([]byte) error {
	*xb = true
	return nil
}
