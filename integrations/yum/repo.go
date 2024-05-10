package yum

import (
	"bytes"
	"compress/gzip"
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/xml"
	"errors"
	"io/fs"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/cavaliergopher/rpm"

	"github.com/kubri/kubri/pkg/crypto/pgp"
	"github.com/kubri/kubri/target"
)

type repo struct {
	primary   *MetaData
	filelists *FileLists
	other     *Other

	dir   string
	files []string
}

func openRepo(ctx context.Context, t target.Target) (*repo, error) {
	dir, err := os.MkdirTemp("", "")
	if err != nil {
		return nil, err
	}

	res := &repo{dir: dir}

	var r RepoMD
	if err := readXML(ctx, t, "repodata/repomd.xml", &r); err != nil {
		if errors.Is(err, fs.ErrNotExist) {
			res.primary = &MetaData{}
			res.filelists = &FileLists{}
			res.other = &Other{}
			return res, nil
		}
		return nil, err
	}

	for _, v := range r.Data {
		var r any
		switch v.Type {
		case "primary":
			r = &res.primary
		case "filelists":
			r = &res.filelists
		case "other":
			r = &res.other
		default:
			continue
		}
		if err := readXML(ctx, t, v.Location.HREF, r); err != nil {
			return nil, err
		}
		res.files = append(res.files, v.Location.HREF)
	}

	if res.primary == nil || res.filelists == nil || res.other == nil {
		return nil, errors.New("invalid repomd.xml")
	}

	return res, nil
}

//nolint:funlen
func (r *repo) Add(b []byte) error {
	h, err := rpm.Read(bytes.NewReader(b))
	if err != nil {
		return err
	}

	id := h.String()
	checksum := sha256.Sum256(b) // TODO: add checksum
	start, end := h.HeaderRange()
	files := h.Files()

	p := Package{
		Type: "rpm",
		Name: h.Name(),
		Arch: h.Architecture(),
		Version: Version{
			Ver:   h.Version(),
			Rel:   h.Release(),
			Epoch: strconv.Itoa(h.Epoch()),
		},
		Checksum: Checksum{
			Type:  "sha256",
			PkgID: "YES",
			Value: hex.EncodeToString(checksum[:]),
		},
		Summary:     h.Summary(),
		Description: h.Description(),
		Packager:    h.Packager(),
		URL:         h.URL(),
		Time: Time{
			File:  timeNow(),
			Build: int(h.BuildTime().Unix()),
		},
		Size: Size{
			Package:   len(b),
			Archive:   int(h.ArchiveSize()),
			Installed: int(h.Size()),
		},
		Location: Location{
			HREF: "Packages/" + id[0:1] + "/" + id + ".rpm",
		},
		Format: Format{
			License:     h.License(),
			Vendor:      h.Vendor(),
			Group:       h.Groups(),
			BuildHost:   h.BuildHost(),
			SourceRPM:   h.SourceRPM(),
			HeaderRange: HeaderRange{Start: start, End: end},
			Provides:    getEntries(h.Provides()),
			Obsoletes:   getEntries(h.Obsoletes()),
			Requires:    getEntries(h.Requires()),
			Conflicts:   getEntries(h.Conflicts()),
			Files:       filterPackageFiles(files),
		},
	}

	r.primary.Package = append(r.primary.Package, p)

	r.filelists.Package = append(r.filelists.Package, FileListsPackage{
		Name:    p.Name,
		PkgID:   p.Checksum.Value,
		Arch:    p.Arch,
		Version: p.Version,
		Files:   convertPackageFiles(files),
	})

	r.other.Package = append(r.other.Package, OtherPackage{
		Name:    p.Name,
		PkgID:   p.Checksum.Value,
		Arch:    p.Arch,
		Version: p.Version,
	})

	return writeFile(filepath.Join(r.dir, p.Location.HREF), b)
}

//nolint:funlen
func (r *repo) Write(pgpKey *pgp.PrivateKey) error {
	md := &RepoMD{}

	data := map[string]any{
		"primary":   r.primary,
		"filelists": r.filelists,
		"other":     r.other,
	}

	for _, name := range []string{"primary", "filelists", "other"} {
		raw, err := xmlMarshal(data[name])
		if err != nil {
			return err
		}

		gz, err := compress(raw)
		if err != nil {
			return err
		}

		var d Data
		d.Type = name
		d.Checksum = getChecksum(gz)
		d.OpenChecksum = getChecksum(raw)
		d.Location.HREF = "repodata/" + d.Checksum.Value + "-" + name + ".xml.gz"
		d.Timestamp = timeNow()
		d.Size = len(gz)
		d.OpenSize = len(raw)

		err = writeFile(filepath.Join(r.dir, d.Location.HREF), gz)
		if err != nil {
			return err
		}

		md.Data = append(md.Data, d)
	}

	md.Revision = timeNow()
	filename := filepath.Join(r.dir, "repodata/repomd.xml")

	b, err := xmlMarshal(md)
	if err != nil {
		return err
	}
	if err = writeFile(filename, b); err != nil {
		return err
	}

	if pgpKey != nil {
		key, err := pgp.MarshalPublicKey(pgp.Public(pgpKey))
		if err != nil {
			return err
		}
		if err = writeFile(filename+".key", key); err != nil {
			return err
		}

		sig, err := pgp.Sign(pgpKey, b)
		if err != nil {
			return err
		}
		if err = writeFile(filename+".asc", sig); err != nil {
			return err
		}
	}

	return nil
}

func readXML(ctx context.Context, t target.Target, path string, res any) error {
	rd, err := t.NewReader(ctx, path)
	if err != nil {
		return err
	}
	defer rd.Close()

	r := rd

	if strings.HasSuffix(path, ".gz") {
		r, err = gzip.NewReader(rd)
		if err != nil {
			return err
		}
		defer r.Close()
	}

	return xml.NewDecoder(r).Decode(res)
}

func filterPackageFiles(files []rpm.FileInfo) []string {
	var res []string
	for _, f := range files {
		if !f.IsDir() && strings.HasPrefix(f.Name(), "/etc/") || strings.Contains(f.Name(), "/bin/") {
			res = append(res, f.Name())
		}
	}
	return res
}

func convertPackageFiles(files []rpm.FileInfo) []File {
	res := make([]File, len(files))
	for i, f := range files {
		var typ string
		if f.IsDir() {
			typ = "dir"
		}
		res[i] = File{Type: typ, Path: f.Name()}
	}
	return res
}

func getEntries(d []rpm.Dependency) *Entries {
	if len(d) == 0 {
		return nil
	}

	entries := make([]Entry, 0, len(d))

	for _, d := range d {
		e := Entry{
			Name: d.Name(),
			Ver:  d.Version(),
			Rel:  d.Release(),
		}
		if e.Ver != "" {
			e.Epoch = strconv.Itoa(d.Epoch())
		}
		entries = append(entries, e)
	}

	return &Entries{entries}
}

func getChecksum(b []byte) Checksum {
	sum := sha256.Sum256(b)
	return Checksum{
		Type:  "sha256",
		Value: hex.EncodeToString(sum[:]),
	}
}

func compress(p []byte) ([]byte, error) {
	var b bytes.Buffer
	w := gzip.NewWriter(&b)
	if _, err := w.Write(p); err != nil {
		return nil, err
	}
	if err := w.Close(); err != nil {
		return nil, err
	}
	return b.Bytes(), nil
}

//nolint:gochecknoglobals
var replacer = strings.NewReplacer("></version>", "/>", "></time>", "/>", "></size>", "/>",
	"></location>", "/>", "></rpm:entry>", "/>", "></rpm:header-range>", "/>")

func xmlMarshal(v any) ([]byte, error) {
	w := bytes.NewBufferString(xml.Header)
	b, err := xml.MarshalIndent(v, "", "\t")
	if err != nil {
		return nil, err
	}
	if _, err = replacer.WriteString(w, string(b)); err != nil {
		return nil, err
	}
	return w.Bytes(), nil
}

func writeFile(path string, data []byte) error {
	if err := os.MkdirAll(filepath.Dir(path), fs.ModePerm); err != nil {
		return err
	}
	return os.WriteFile(path, data, 0o600)
}

//nolint:gochecknoglobals
var timeNow = func() int { return int(time.Now().Unix()) }
