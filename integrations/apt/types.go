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

	b := make([]byte, 0, 512)
	for _, info := range c {
		b = append(b, fmt.Sprintf("\n%x %d %s", info.Sum, info.Size, info.Name)...)
	}

	return b, nil
}

var space = []byte{' '} //nolint:gochecknoglobals

func (c *Checksums[T]) UnmarshalText(b []byte) error {
	r := bufio.NewScanner(bytes.NewReader(b))
	if !r.Scan() {
		return nil
	}

	for r.Scan() {
		a := bytes.SplitN(r.Bytes(), space, 3)
		size, err := strconv.Atoi(string(a[1]))
		if err != nil {
			return err
		}

		var sum T
		if hex.DecodedLen(len(a[0])) != len(sum) {
			return errors.New("hex data would overflow byte array")
		}
		b := make([]byte, len(sum))
		if _, err = hex.Decode(b, a[0]); err != nil {
			return err
		}
		for i := 0; i < len(sum); i++ {
			sum[i] = b[i]
		}

		*c = append(*c, ChecksumFile[T]{sum, size, string(a[2])})
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
