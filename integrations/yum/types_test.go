package yum_test

import (
	"encoding/xml"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"

	"github.com/kubri/kubri/integrations/yum"
)

func TestRepoMDMarshalUnmarshal(t *testing.T) {
	in := &yum.RepoMD{
		Data: []yum.Data{
			{
				Type: "primary",
				Location: yum.Location{
					HREF: "repodata/primary.xml.gz",
				},
				Checksum: yum.Checksum{
					Type:  "sha256",
					Value: "abc",
				},
				OpenChecksum: yum.Checksum{
					Type:  "sha256",
					Value: "abc",
				},
				Timestamp: 1192033773,
			},
			{
				Type: "filelists",
				Location: yum.Location{
					HREF: "repodata/filelists.xml.gz",
				},
				Checksum: yum.Checksum{
					Type:  "sha256",
					Value: "abc",
				},
				OpenChecksum: yum.Checksum{
					Type:  "sha256",
					Value: "abc",
				},
				Timestamp: 1192033773,
			},
			{
				Type: "other",
				Location: yum.Location{
					HREF: "repodata/other.xml.gz",
				},
				Checksum: yum.Checksum{
					Type:  "sha256",
					Value: "abc",
				},
				OpenChecksum: yum.Checksum{
					Type:  "sha256",
					Value: "abc",
				},
				Timestamp: 1192033773,
			},
		},
	}

	want := `<repomd xmlns="http://linux.duke.edu/metadata/repo">
	<data type="primary">
		<checksum type="sha256">abc</checksum>
		<open-checksum type="sha256">abc</open-checksum>
		<location href="repodata/primary.xml.gz"></location>
		<timestamp>1192033773</timestamp>
	</data>
	<data type="filelists">
		<checksum type="sha256">abc</checksum>
		<open-checksum type="sha256">abc</open-checksum>
		<location href="repodata/filelists.xml.gz"></location>
		<timestamp>1192033773</timestamp>
	</data>
	<data type="other">
		<checksum type="sha256">abc</checksum>
		<open-checksum type="sha256">abc</open-checksum>
		<location href="repodata/other.xml.gz"></location>
		<timestamp>1192033773</timestamp>
	</data>
</repomd>`

	testMarshalUnmarshal(t, in, want)
}

func TestMetaDataMarshalUnmarshal(t *testing.T) {
	in := &yum.MetaData{
		Package: []yum.Package{{
			Type: "rpm",
			Name: "Title",
			Arch: "x86_64",
			Version: yum.Version{
				Ver:   "1.0",
				Rel:   "1",
				Epoch: "0",
			},
			Checksum: yum.Checksum{
				Type:  "sha256",
				Value: "abc",
			},
			Summary:     "Summary",
			Description: "Description",
			Time: yum.Time{
				File:  1192033774,
				Build: 1192033773,
			},
			Size: yum.Size{
				Package:   1000,
				Installed: 2000,
				Archive:   2000,
			},
			Location: yum.Location{
				HREF: "Packages/test.x86_64.rpm",
			},
			Format: yum.Format{
				License: "MIT",
				HeaderRange: yum.HeaderRange{
					Start: 5,
					End:   10,
				},
			},
		}},
	}

	want := `<metadata xmlns="http://linux.duke.edu/metadata/common" xmlns:rpm="http://linux.duke.edu/metadata/rpm" packages="1">
	<package type="rpm">
		<name>Title</name>
		<arch>x86_64</arch>
		<version ver="1.0" rel="1" epoch="0"></version>
		<checksum type="sha256">abc</checksum>
		<summary>Summary</summary>
		<description>Description</description>
		<time file="1192033774" build="1192033773"></time>
		<size package="1000" installed="2000" archive="2000"></size>
		<location href="Packages/test.x86_64.rpm"></location>
		<format>
			<rpm:license>MIT</rpm:license>
			<rpm:header-range start="5" end="10"></rpm:header-range>
		</format>
	</package>
</metadata>`

	testMarshalUnmarshal(t, in, want)
}

func testMarshalUnmarshal[T any](t *testing.T, in *T, want string) {
	t.Helper()

	b, err := xml.MarshalIndent(in, "", "\t")
	if err != nil {
		t.Fatal(err)
	}

	if diff := cmp.Diff(want, string(b)); diff != "" {
		t.Fatal(diff)
	}

	got := new(T)
	if err = xml.Unmarshal(b, got); err != nil {
		t.Fatal(err)
	}

	if diff := cmp.Diff(in, got, cmpopts.IgnoreTypes(xml.Name{})); diff != "" {
		t.Fatal(diff)
	}
}
