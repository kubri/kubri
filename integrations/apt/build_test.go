package apt_test

import (
	"context"
	"io"
	"path"
	"testing"
	"time"

	"github.com/abemedia/appcast/integrations/apt"
	"github.com/abemedia/appcast/pkg/crypto/pgp"
	source "github.com/abemedia/appcast/source/file"
	"github.com/abemedia/appcast/target"
	ftarget "github.com/abemedia/appcast/target/file"
	"github.com/google/go-cmp/cmp"
	"github.com/klauspost/compress/gzip"
)

func TestBuild(t *testing.T) {
	now := time.Now()

	want := map[string]string{
		"dists/stable/InRelease": "",
		"dists/stable/Release": `Suite: stable
Codename: stable
Date: ` + now.UTC().Format(time.RFC1123) + `
Architectures: amd64 i386
Components: main
MD5Sum:
 caa0d53d538fb71d5421947aa7cbdf62 440 main/binary-amd64/Packages
 e7cb8e6ad0c014da8f9138bb0e22e211 344 main/binary-amd64/Packages.gz
 ed50c229b0405e82b42df9dffd1a478e 66 main/binary-amd64/Release
 7b8f49611c8a35fcdb806334c93e34b4 438 main/binary-i386/Packages
 dd8a037bfe8eac16164fe3c6427e393a 344 main/binary-i386/Packages.gz
 1c76c0e691587e8c0fa470f96791a23f 65 main/binary-i386/Release
SHA1:
 88414e55f369a13827e06e684f8930da93d5baf7 440 main/binary-amd64/Packages
 af6f87cf3df086a6f12aa36908021dd8470be55b 344 main/binary-amd64/Packages.gz
 0356f9e4028428574fb61279d6b7e3f64e02c474 66 main/binary-amd64/Release
 06a944d8704b911608bd9923ceaf0d309fcbb1ad 438 main/binary-i386/Packages
 70b38c9b46e4495f73941e8b25b283c39935f4ba 344 main/binary-i386/Packages.gz
 4c939f94794ad0f1ab92dedf4a3f8b9323a8730d 65 main/binary-i386/Release
SHA256:
 38fc8aed25c62a09260caef1b34c4411f1c862193cefe05c60e36b8a49e6a97e 440 main/binary-amd64/Packages
 9b2875e549ce8e73501002362c9ef84ce7170a6a1263888316194f041a2b6d8c 344 main/binary-amd64/Packages.gz
 eead410c09874826a4e597f34ae597d9e1eda87b072ebad5983d835c15139fe8 66 main/binary-amd64/Release
 321277b81510da1f31fd95090c8ffcc75918a8d364c4d5fbbe2fd6b1dd29228b 438 main/binary-i386/Packages
 312116a1da18b05c67d1bf4f6987e7998188382a8e36999eb8840020454add4a 344 main/binary-i386/Packages.gz
 133ff0295c5a0d85673cd4d65a23ac847f62c0ffd68b706d631e61f4465f3b72 65 main/binary-i386/Release
`,
		"dists/stable/main/binary-amd64/Packages": `Package: test
Version: 1.0
Architecture: amd64
Maintainer: Test User <test@example.com>
Priority: optional
Section: utils
Filename: pool/main/t/test/test_1.0_amd64.deb
Size: 614
MD5sum: ac900eaebfdb5081dd5c0138ccd8d652
SHA1: 55440d5e25c2550d4f0f60a8b743e857bb022fc4
SHA256: 88b64d49cb11b27af379f16a6bdc7dd6da34927c470a76e7199fd2d04eee204c
Homepage: https://example.com
Description: This is a test.
 It does nothing.
 .
 Absolutely nothing.
`,
		"dists/stable/main/binary-amd64/Packages.gz": "",
		"dists/stable/main/binary-amd64/Release":     "Archive: stable\nSuite: stable\nArchitecture: amd64\nComponent: main\n",
		"dists/stable/main/binary-i386/Packages": `Package: test
Version: 1.0
Architecture: i386
Maintainer: Test User <test@example.com>
Priority: optional
Section: utils
Filename: pool/main/t/test/test_1.0_i386.deb
Size: 610
MD5sum: 8db9c79bd45747682792b9a968a973c6
SHA1: 63e4f301d0094df113d6b4a737fb48af6f742611
SHA256: a71afe44cdd296db564c1de7f36026f96be964afe85903288627ec2582638b25
Homepage: https://example.com
Description: This is a test.
 It does nothing.
 .
 Absolutely nothing.
`,
		"dists/stable/main/binary-i386/Packages.gz": "",
		"dists/stable/main/binary-i386/Release":     "Archive: stable\nSuite: stable\nArchitecture: i386\nComponent: main\n",
	}

	c := &apt.Config{}
	c.Source, _ = source.New(source.Config{Path: "../../testdata"})
	c.Target, _ = ftarget.New(ftarget.Config{Path: t.TempDir()})

	testBuild(t, c, want, now)

	// Should be no-op as nothing changed so timestamp should still be valid.
	testBuild(t, c, want, now.Add(time.Hour))

	t.Run("PGP", func(t *testing.T) {
		c.Target, _ = ftarget.New(ftarget.Config{Path: t.TempDir()})
		c.PGPKey, _ = pgp.NewPrivateKey("test", "test@example.com")
		want["dists/stable/Release.gpg"] = ""
		testBuild(t, c, want, now)
	})
}

func testBuild(t *testing.T, c *apt.Config, want map[string]string, now time.Time) { //nolint:thelper
	apt.SetTime(now)

	err := apt.Build(context.Background(), c)
	if err != nil {
		t.Fatal(err)
	}

	for name, data := range want {
		got := readFile(t, c.Target, name)

		switch path.Base(name) {
		case "Packages.gz":
			data = want[name[:len(name)-3]]
		case "InRelease":
			data = want[path.Dir(name)+"/Release"]
			if c.PGPKey != nil {
				var sig []byte
				got, sig, err = pgp.Split(got)
				if err != nil {
					t.Fatal(err)
				}
				if !pgp.Verify(pgp.Public(c.PGPKey), got, sig) {
					t.Error(name, "failed to verify signature")
				}
			}
		case "Release.gpg":
			if !pgp.Verify(pgp.Public(c.PGPKey), readFile(t, c.Target, path.Dir(name)+"/Release"), got) {
				t.Error(name, "failed to verify signature")
			}
			continue
		}

		if diff := cmp.Diff(data, string(got)); diff != "" {
			t.Error(name, diff)
		}
	}
}

func readFile(t *testing.T, tgt target.Target, name string) []byte {
	t.Helper()

	r, err := tgt.NewReader(context.Background(), name)
	if err != nil {
		t.Fatal(name, err)
	}
	defer r.Close()

	if path.Ext(name) == ".gz" {
		r, err = gzip.NewReader(r)
		if err != nil {
			t.Fatal(name, err)
		}
	}

	b, err := io.ReadAll(r)
	if err != nil {
		t.Fatal(name, err)
	}

	return b
}
