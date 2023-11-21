package yum

import (
	"encoding/xml"
	"strconv"
)

type RepoMD struct {
	XMLName  xml.Name `xml:"http://linux.duke.edu/metadata/repo repomd"`
	Revision int      `xml:"revision,omitempty"`
	Data     []Data   `xml:"data"`
}

type Data struct {
	Type         string   `xml:"type,attr"`
	Checksum     Checksum `xml:"checksum"`
	OpenChecksum Checksum `xml:"open-checksum"`
	Location     Location `xml:"location"`
	Timestamp    int      `xml:"timestamp"`
	Size         int      `xml:"size,omitempty"`
	OpenSize     int      `xml:"open-size,omitempty"`
}

type Checksum struct {
	Type  string `xml:"type,attr"`
	PkgID string `xml:"pkgid,attr,omitempty"`
	Value string `xml:",chardata"`
}

type Location struct {
	HREF string `xml:"href,attr"`
}

type MetaData struct {
	Package []Package `xml:"package"`
}

func (r *MetaData) MarshalXML(enc *xml.Encoder, _ xml.StartElement) error {
	err := enc.EncodeToken(xml.StartElement{
		Name: xml.Name{Local: "metadata"},
		Attr: []xml.Attr{
			{Name: xml.Name{Local: "xmlns"}, Value: "http://linux.duke.edu/metadata/common"},
			{Name: xml.Name{Local: "xmlns:rpm"}, Value: "http://linux.duke.edu/metadata/rpm"},
			{Name: xml.Name{Local: "packages"}, Value: strconv.Itoa(len(r.Package))},
		},
	})
	if err != nil {
		return err
	}

	err = enc.EncodeElement(r.Package, xml.StartElement{Name: xml.Name{Local: "package"}})
	if err != nil {
		return err
	}

	return enc.EncodeToken(xml.EndElement{Name: xml.Name{Local: "metadata"}})
}

type Package struct {
	Type        string   `xml:"type,attr"`
	Name        string   `xml:"name"`
	Arch        string   `xml:"arch"`
	Version     Version  `xml:"version"`
	Checksum    Checksum `xml:"checksum"`
	Summary     string   `xml:"summary"`
	Description string   `xml:"description"`
	Packager    string   `xml:"packager,omitempty"`
	URL         string   `xml:"url,omitempty"`
	Time        Time     `xml:"time"`
	Size        Size     `xml:"size"`
	Location    Location `xml:"location"`
	Format      Format   `xml:"format"`
}

type Time struct {
	File  int `xml:"file,attr"`
	Build int `xml:"build,attr"`
}

type Format struct {
	License     string      `xml:"rpm:license,omitempty"`
	Vendor      string      `xml:"rpm:vendor,omitempty"`
	Group       string      `xml:"rpm:group,omitempty"`
	BuildHost   string      `xml:"rpm:buildhost,omitempty"`
	SourceRPM   string      `xml:"rpm:sourcerpm,omitempty"`
	HeaderRange HeaderRange `xml:"rpm:header-range"`
	Provides    *Entries    `xml:"rpm:provides,omitempty"`
	Obsoletes   *Entries    `xml:"rpm:obsoletes,omitempty"`
	Requires    *Entries    `xml:"rpm:requires,omitempty"`
	Conflicts   *Entries    `xml:"rpm:conflicts,omitempty"`
	Files       []string    `xml:"file,omitempty"`
}

func (f *Format) UnmarshalXML(dec *xml.Decoder, start xml.StartElement) error {
	type entries struct {
		Entries []Entry `xml:"entry"`
	}
	var data struct {
		License     string      `xml:"license"`
		Vendor      string      `xml:"vendor"`
		Group       string      `xml:"group"`
		BuildHost   string      `xml:"buildhost"`
		SourceRPM   string      `xml:"sourcerpm"`
		HeaderRange HeaderRange `xml:"header-range"`
		Provides    *entries    `xml:"provides"`
		Obsoletes   *entries    `xml:"obsoletes"`
		Requires    *entries    `xml:"requires"`
		Conflicts   *entries    `xml:"conflicts"`
		Files       []string    `xml:"file"`
	}
	if err := dec.DecodeElement(&data, &start); err != nil {
		return err
	}

	*f = Format{
		License:     data.License,
		Vendor:      data.Vendor,
		Group:       data.Group,
		BuildHost:   data.BuildHost,
		SourceRPM:   data.SourceRPM,
		HeaderRange: data.HeaderRange,
		Provides:    (*Entries)(data.Provides),
		Obsoletes:   (*Entries)(data.Obsoletes),
		Requires:    (*Entries)(data.Requires),
		Conflicts:   (*Entries)(data.Conflicts),
		Files:       data.Files,
	}

	return nil
}

type HeaderRange struct {
	Start int `xml:"start,attr"`
	End   int `xml:"end,attr"`
}

type Entries struct {
	Entries []Entry `xml:"rpm:entry,omitempty"`
}

type Entry struct {
	Name  string `xml:"name,attr"`
	Flags string `xml:"flags,attr,omitempty"`
	Epoch string `xml:"epoch,attr,omitempty"`
	Ver   string `xml:"ver,attr,omitempty"`
	Rel   string `xml:"rel,attr,omitempty"`
	Pre   string `xml:"pre,attr,omitempty"`
}

type Size struct {
	Package   int `xml:"package,attr"`
	Installed int `xml:"installed,attr"`
	Archive   int `xml:"archive,attr"`
}

type FileLists struct {
	Packages string             `xml:"packages,attr"`
	Package  []FileListsPackage `xml:"package"`
}

func (r *FileLists) MarshalXML(enc *xml.Encoder, _ xml.StartElement) error {
	err := enc.EncodeToken(xml.StartElement{
		Name: xml.Name{Local: "filelists"},
		Attr: []xml.Attr{
			{Name: xml.Name{Local: "xmlns"}, Value: "http://linux.duke.edu/metadata/filelists"},
			{Name: xml.Name{Local: "packages"}, Value: strconv.Itoa(len(r.Package))},
		},
	})
	if err != nil {
		return err
	}

	err = enc.EncodeElement(r.Package, xml.StartElement{Name: xml.Name{Local: "package"}})
	if err != nil {
		return err
	}

	return enc.EncodeToken(xml.EndElement{Name: xml.Name{Local: "filelists"}})
}

type FileListsPackage struct {
	Name    string  `xml:"name,attr"`
	PkgID   string  `xml:"pkgid,attr"`
	Arch    string  `xml:"arch,attr"`
	Version Version `xml:"version"`
	Files   []File  `xml:"file"`
}

type File struct {
	Type string `xml:"type,attr,omitempty"`
	Path string `xml:",chardata"`
}

type Version struct {
	Ver   string `xml:"ver,attr"`
	Rel   string `xml:"rel,attr"`
	Epoch string `xml:"epoch,attr"`
}

type Other struct {
	Package []OtherPackage `xml:"package"`
}

func (r *Other) MarshalXML(enc *xml.Encoder, _ xml.StartElement) error {
	err := enc.EncodeToken(xml.StartElement{
		Name: xml.Name{Local: "other"},
		Attr: []xml.Attr{
			{Name: xml.Name{Local: "xmlns"}, Value: "http://linux.duke.edu/metadata/other"},
			{Name: xml.Name{Local: "packages"}, Value: strconv.Itoa(len(r.Package))},
		},
	})
	if err != nil {
		return err
	}

	err = enc.EncodeElement(r.Package, xml.StartElement{Name: xml.Name{Local: "package"}})
	if err != nil {
		return err
	}

	return enc.EncodeToken(xml.EndElement{Name: xml.Name{Local: "other"}})
}

type OtherPackage struct {
	Name    string  `xml:"name,attr"`
	PkgID   string  `xml:"pkgid,attr"`
	Arch    string  `xml:"arch,attr"`
	Version Version `xml:"version"`
}
