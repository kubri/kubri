package apt

import (
	"bufio"
	"bytes"
	"encoding/hex"
	"errors"
	"fmt"
	"strconv"
	"time"
)

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
	MD5Sum        Checksums[[16]byte]
	SHA1          Checksums[[20]byte]
	SHA256        Checksums[[32]byte]
	SHA512        Checksums[[64]byte]
}

type sum interface {
	[16]byte | [20]byte | [32]byte | [64]byte
}

type Checksums[T sum] []ChecksumFile[T]

func (c Checksums[T]) MarshalText() ([]byte, error) {
	if len(c) == 0 {
		return nil, nil
	}

	b := make([]byte, 0, len(c)*(hex.EncodedLen(len(c[0].Sum))+40))
	for _, info := range c {
		b = fmt.Appendf(b, "\n%x %d %s", info.Sum, info.Size, info.Name)
	}

	return b, nil
}

var space = []byte{' '} //nolint:gochecknoglobals

func (c *Checksums[T]) UnmarshalText(b []byte) (err error) {
	r := bufio.NewScanner(bytes.NewReader(b))
	if !r.Scan() {
		return nil
	}

	var f ChecksumFile[T]

	for r.Scan() {
		sum, rest, _ := bytes.Cut(r.Bytes(), space)
		size, name, _ := bytes.Cut(rest, space)

		f.Name = string(name)
		f.Size, err = strconv.Atoi(string(size))
		if err != nil {
			return err
		}

		if hex.DecodedLen(len(sum)) != len(f.Sum) {
			return errors.New("hex data would overflow byte array")
		}
		if _, err = hex.Decode(sum, sum); err != nil {
			return err
		}
		for i := 0; i < len(f.Sum); i++ {
			f.Sum[i] = sum[i]
		}

		*c = append(*c, f)
	}

	return r.Err()
}

type ChecksumFile[T sum] struct {
	Sum  T
	Size int
	Name string
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
