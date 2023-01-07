package apt

import "time"

type Release struct {
	Origin       string
	Label        string
	Archive      string
	Suite        string
	Architecture string
	Component    string
	Description  string
}

type Releases struct {
	Origin        string
	Label         string
	Suite         string
	Codename      string
	Date          time.Time
	Architectures string
	Components    string
	Description   string
	MD5Sum        string
	SHA1          string
	SHA256        string
	SHA512        string
}

type Package struct {
	Package       string
	Version       string
	Architecture  string
	Maintainer    string
	InstalledSize int    `deb:"Installed-Size"`
	PreDepends    string `deb:"Pre-Depends"`
	Depends       string
	Recommends    string
	Conflicts     string
	Replaces      string
	Provides      string
	Priority      string
	Section       string
	Filename      string
	Size          int
	MD5sum        [16]byte
	SHA1          [20]byte
	SHA256        [32]byte
	Homepage      string
	Description   string
}
