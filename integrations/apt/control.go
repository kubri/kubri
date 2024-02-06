package apt

import (
	"archive/tar"
	"bytes"
	"io"
	"path"
	"strings"

	"github.com/blakesmith/ar"

	"github.com/kubri/kubri/integrations/apt/deb"
)

func getControl(b []byte) (*Package, error) {
	r := ar.NewReader(bytes.NewReader(b))

	for {
		h, err := r.Next()
		if err != nil {
			return nil, err
		}
		if !strings.HasPrefix(h.Name, "control.tar") {
			continue
		}

		r, err := decompress(path.Ext(h.Name))(r)
		if err != nil {
			return nil, err
		}
		defer r.Close()

		tr := tar.NewReader(r)

		for {
			h, err := tr.Next()
			if err != nil {
				return nil, err
			}
			if path.Base(h.Name) != "control" {
				continue
			}

			b, err := io.ReadAll(tr)
			if err != nil {
				return nil, err
			}
			p := &Package{}
			if err = deb.Unmarshal(b, p); err != nil {
				return nil, err
			}

			return p, nil
		}
	}
}
