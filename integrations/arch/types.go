package arch

// Package holds minimal metadata for a package.
type Package struct {
	Filename       string   `desc:"FILENAME"    pkginfo:"-"`
	Name           string   `desc:"NAME"        pkginfo:"pkgname"`
	Base           string   `desc:"BASE"        pkginfo:"pkgbase"`
	Version        string   `desc:"VERSION"     pkginfo:"pkgver"`
	Desc           string   `desc:"DESC"        pkginfo:"pkgdesc"`
	CompressedSize int64    `desc:"CSIZE"       pkginfo:"-"`
	InstalledSize  int64    `desc:"ISIZE"       pkginfo:"size"`
	SHA256Sum      [32]byte `desc:"SHA256SUM"   pkginfo:"-"`
	URL            string   `desc:"URL"         pkginfo:"url"`
	License        string   `desc:"LICENSE"     pkginfo:"license"`
	Arch           string   `desc:"ARCH"        pkginfo:"arch"`
	BuildDate      string   `desc:"BUILDDATE"   pkginfo:"builddate"`
	Packager       string   `desc:"PACKAGER"    pkginfo:"packager"`
	Replaces       []string `desc:"REPLACES"    pkginfo:"replaces"`
	Conflicts      []string `desc:"CONFLICTS"   pkginfo:"conflict"`
	Provides       []string `desc:"PROVIDES"    pkginfo:"provides"`
	Depends        []string `desc:"DEPENDS"     pkginfo:"depend"`
	OptDepends     []string `desc:"OPTDEPENDS"  pkginfo:"optdepend"`
	MakeDepends    []string `desc:"MAKEDEPENDS" pkginfo:"makedepend"`
}

// Files holds the list of files in a package.
type Files struct {
	Files string `desc:"FILES" pkginfo:"-"`
}
