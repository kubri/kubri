package apt_test

import (
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"

	"github.com/kubri/kubri/integrations/apt"
	"github.com/kubri/kubri/integrations/apt/deb"
)

func TestReleases(t *testing.T) {
	in := &apt.Releases{
		Suite:         "stable",
		Codename:      "stable",
		Date:          time.Date(2023, 11, 19, 16, 30, 23, 0, time.UTC),
		Architectures: "amd64 i386",
		Components:    "main",
		MD5Sum: apt.Checksums[[16]byte]{
			{
				Sum:  [16]byte{202, 160, 213, 61, 83, 143, 183, 29, 84, 33, 148, 122, 167, 203, 223, 98},
				Size: 440,
				Name: "main/binary-amd64/Packages",
			},
			{
				Sum:  [16]byte{231, 203, 142, 106, 208, 192, 20, 218, 143, 145, 56, 187, 14, 34, 226, 17},
				Size: 344,
				Name: "main/binary-amd64/Packages.gz",
			},
			{
				Sum:  [16]byte{96, 119, 224, 33, 219, 233, 58, 68, 227, 141, 92, 99, 227, 135, 96, 113},
				Size: 412,
				Name: "main/binary-amd64/Packages.xz",
			},
			{
				Sum:  [16]byte{237, 80, 194, 41, 176, 64, 94, 130, 180, 45, 249, 223, 253, 26, 71, 142},
				Size: 66,
				Name: "main/binary-amd64/Release",
			},
			{
				Sum:  [16]byte{123, 143, 73, 97, 28, 138, 53, 252, 219, 128, 99, 52, 201, 62, 52, 180},
				Size: 438,
				Name: "main/binary-i386/Packages",
			},
			{
				Sum:  [16]byte{221, 138, 3, 123, 254, 142, 172, 22, 22, 79, 227, 198, 66, 126, 57, 58},
				Size: 344,
				Name: "main/binary-i386/Packages.gz",
			},
			{
				Sum:  [16]byte{24, 245, 32, 137, 124, 185, 111, 29, 135, 235, 190, 128, 29, 11, 33, 74},
				Size: 412,
				Name: "main/binary-i386/Packages.xz",
			},
			{
				Sum:  [16]byte{28, 118, 192, 230, 145, 88, 126, 140, 15, 164, 112, 249, 103, 145, 162, 63},
				Size: 65,
				Name: "main/binary-i386/Release",
			},
		},
		SHA1: apt.Checksums[[20]byte]{
			{
				Sum:  [20]byte{136, 65, 78, 85, 243, 105, 161, 56, 39, 224, 110, 104, 79, 137, 48, 218, 147, 213, 186, 247},
				Size: 440,
				Name: "main/binary-amd64/Packages",
			},
			{
				Sum:  [20]byte{175, 111, 135, 207, 61, 240, 134, 166, 241, 42, 163, 105, 8, 2, 29, 216, 71, 11, 229, 91},
				Size: 344,
				Name: "main/binary-amd64/Packages.gz",
			},
			{
				Sum:  [20]byte{196, 127, 66, 82, 175, 57, 212, 194, 87, 86, 224, 211, 218, 156, 11, 8, 99, 6, 240, 241},
				Size: 412,
				Name: "main/binary-amd64/Packages.xz",
			},
			{
				Sum:  [20]byte{3, 86, 249, 228, 2, 132, 40, 87, 79, 182, 18, 121, 214, 183, 227, 246, 78, 2, 196, 116},
				Size: 66,
				Name: "main/binary-amd64/Release",
			},
			{
				Sum:  [20]byte{6, 169, 68, 216, 112, 75, 145, 22, 8, 189, 153, 35, 206, 175, 13, 48, 159, 203, 177, 173},
				Size: 438,
				Name: "main/binary-i386/Packages",
			},
			{
				Sum:  [20]byte{112, 179, 140, 155, 70, 228, 73, 95, 115, 148, 30, 139, 37, 178, 131, 195, 153, 53, 244, 186},
				Size: 344,
				Name: "main/binary-i386/Packages.gz",
			},
			{
				Sum:  [20]byte{38, 10, 245, 60, 182, 128, 9, 229, 80, 25, 1, 51, 216, 140, 0, 224, 91, 72, 167, 43},
				Size: 412,
				Name: "main/binary-i386/Packages.xz",
			},
			{
				Sum:  [20]byte{76, 147, 159, 148, 121, 74, 208, 241, 171, 146, 222, 223, 74, 63, 139, 147, 35, 168, 115, 13},
				Size: 65,
				Name: "main/binary-i386/Release",
			},
		},
		SHA256: apt.Checksums[[32]byte]{
			{
				Sum:  [32]byte{56, 252, 138, 237, 37, 198, 42, 9, 38, 12, 174, 241, 179, 76, 68, 17, 241, 200, 98, 25, 60, 239, 224, 92, 96, 227, 107, 138, 73, 230, 169, 126},
				Size: 440,
				Name: "main/binary-amd64/Packages",
			},
			{
				Sum:  [32]byte{155, 40, 117, 229, 73, 206, 142, 115, 80, 16, 2, 54, 44, 158, 248, 76, 231, 23, 10, 106, 18, 99, 136, 131, 22, 25, 79, 4, 26, 43, 109, 140},
				Size: 344,
				Name: "main/binary-amd64/Packages.gz",
			},
			{
				Sum:  [32]byte{155, 254, 54, 168, 31, 228, 124, 215, 80, 249, 57, 93, 59, 240, 215, 107, 247, 183, 76, 247, 255, 72, 122, 184, 197, 247, 85, 47, 194, 88, 103, 161},
				Size: 412,
				Name: "main/binary-amd64/Packages.xz",
			},
			{
				Sum:  [32]byte{238, 173, 65, 12, 9, 135, 72, 38, 164, 229, 151, 243, 74, 229, 151, 217, 225, 237, 168, 123, 7, 46, 186, 213, 152, 61, 131, 92, 21, 19, 159, 232},
				Size: 66,
				Name: "main/binary-amd64/Release",
			},
			{
				Sum:  [32]byte{50, 18, 119, 184, 21, 16, 218, 31, 49, 253, 149, 9, 12, 143, 252, 199, 89, 24, 168, 211, 100, 196, 213, 251, 190, 47, 214, 177, 221, 41, 34, 139},
				Size: 438,
				Name: "main/binary-i386/Packages",
			},
			{
				Sum:  [32]byte{49, 33, 22, 161, 218, 24, 176, 92, 103, 209, 191, 79, 105, 135, 231, 153, 129, 136, 56, 42, 142, 54, 153, 158, 184, 132, 0, 32, 69, 74, 221, 74},
				Size: 344,
				Name: "main/binary-i386/Packages.gz",
			},
			{
				Sum:  [32]byte{99, 126, 21, 136, 122, 255, 22, 74, 69, 123, 119, 32, 30, 84, 135, 118, 140, 224, 233, 14, 174, 125, 148, 239, 72, 48, 99, 128, 109, 193, 100, 163},
				Size: 412,
				Name: "main/binary-i386/Packages.xz",
			},
			{
				Sum:  [32]byte{19, 63, 240, 41, 92, 90, 13, 133, 103, 60, 212, 214, 90, 35, 172, 132, 127, 98, 192, 255, 214, 139, 112, 109, 99, 30, 97, 244, 70, 95, 59, 114},
				Size: 65,
				Name: "main/binary-i386/Release",
			},
		},
	}

	want := `Suite: stable
Codename: stable
Date: Sun, 19 Nov 2023 16:30:23 UTC
Architectures: amd64 i386
Components: main
MD5Sum:
 caa0d53d538fb71d5421947aa7cbdf62 440 main/binary-amd64/Packages
 e7cb8e6ad0c014da8f9138bb0e22e211 344 main/binary-amd64/Packages.gz
 6077e021dbe93a44e38d5c63e3876071 412 main/binary-amd64/Packages.xz
 ed50c229b0405e82b42df9dffd1a478e 66 main/binary-amd64/Release
 7b8f49611c8a35fcdb806334c93e34b4 438 main/binary-i386/Packages
 dd8a037bfe8eac16164fe3c6427e393a 344 main/binary-i386/Packages.gz
 18f520897cb96f1d87ebbe801d0b214a 412 main/binary-i386/Packages.xz
 1c76c0e691587e8c0fa470f96791a23f 65 main/binary-i386/Release
SHA1:
 88414e55f369a13827e06e684f8930da93d5baf7 440 main/binary-amd64/Packages
 af6f87cf3df086a6f12aa36908021dd8470be55b 344 main/binary-amd64/Packages.gz
 c47f4252af39d4c25756e0d3da9c0b086306f0f1 412 main/binary-amd64/Packages.xz
 0356f9e4028428574fb61279d6b7e3f64e02c474 66 main/binary-amd64/Release
 06a944d8704b911608bd9923ceaf0d309fcbb1ad 438 main/binary-i386/Packages
 70b38c9b46e4495f73941e8b25b283c39935f4ba 344 main/binary-i386/Packages.gz
 260af53cb68009e550190133d88c00e05b48a72b 412 main/binary-i386/Packages.xz
 4c939f94794ad0f1ab92dedf4a3f8b9323a8730d 65 main/binary-i386/Release
SHA256:
 38fc8aed25c62a09260caef1b34c4411f1c862193cefe05c60e36b8a49e6a97e 440 main/binary-amd64/Packages
 9b2875e549ce8e73501002362c9ef84ce7170a6a1263888316194f041a2b6d8c 344 main/binary-amd64/Packages.gz
 9bfe36a81fe47cd750f9395d3bf0d76bf7b74cf7ff487ab8c5f7552fc25867a1 412 main/binary-amd64/Packages.xz
 eead410c09874826a4e597f34ae597d9e1eda87b072ebad5983d835c15139fe8 66 main/binary-amd64/Release
 321277b81510da1f31fd95090c8ffcc75918a8d364c4d5fbbe2fd6b1dd29228b 438 main/binary-i386/Packages
 312116a1da18b05c67d1bf4f6987e7998188382a8e36999eb8840020454add4a 344 main/binary-i386/Packages.gz
 637e15887aff164a457b77201e5487768ce0e90eae7d94ef483063806dc164a3 412 main/binary-i386/Packages.xz
 133ff0295c5a0d85673cd4d65a23ac847f62c0ffd68b706d631e61f4465f3b72 65 main/binary-i386/Release
`

	testMarshalUnmarshal(t, in, want)
}

func TestRelease(t *testing.T) {
	in := &apt.Release{
		Archive:      "stable",
		Suite:        "stable",
		Architecture: "amd64",
		Component:    "main",
	}

	want := "Archive: stable\nSuite: stable\nArchitecture: amd64\nComponent: main\n"

	testMarshalUnmarshal(t, in, want)
}

func TestPackages(t *testing.T) {
	in := []*apt.Package{
		{
			Package:      "appcast-test",
			Version:      "1.1.0",
			Architecture: "amd64",
			Maintainer:   "Test User <test@example.com>",
			Depends:      "bash",
			Recommends:   "git",
			Conflicts:    "appcast-test-new",
			Replaces:     "appcast-test-old",
			Provides:     "appcast-test-alt",
			Priority:     "optional",
			Section:      "utils",
			Filename:     "pool/main/a/appcast-test/appcast-test_1.1.0_amd64.deb",
			Size:         812,
			MD5sum:       [16]byte{2, 7, 176, 158, 32, 192, 123, 213, 242, 176, 93, 153, 167, 147, 25, 25},
			SHA1:         [20]byte{13, 52, 171, 147, 58, 249, 90, 133, 116, 199, 169, 206, 18, 218, 112, 78, 15, 63, 66, 235},
			SHA256:       [32]byte{249, 182, 206, 160, 145, 46, 86, 254, 34, 183, 190, 244, 163, 136, 151, 116, 58, 87, 237, 159, 232, 78, 1, 27, 179, 117, 240, 51, 190, 102, 55, 246},
			Homepage:     "http://example.com",
			Description:  "This is a test.\nIt does nothing.\n\nAbsolutely nothing.",
		},
		{
			Package:      "appcast-test",
			Version:      "1.0.0",
			Architecture: "amd64",
			Maintainer:   "Test User <test@example.com>",
			Depends:      "bash",
			Recommends:   "git",
			Conflicts:    "appcast-test-new",
			Replaces:     "appcast-test-old",
			Provides:     "appcast-test-alt",
			Priority:     "optional",
			Section:      "utils",
			Filename:     "pool/main/a/appcast-test/appcast-test_1.0.0_amd64.deb",
			Size:         814,
			MD5sum:       [16]byte{28, 21, 209, 214, 104, 4, 205, 7, 56, 105, 232, 188, 34, 60, 245, 41},
			SHA1:         [20]byte{16, 76, 26, 210, 207, 216, 40, 72, 14, 117, 148, 112, 5, 218, 252, 98, 87, 207, 5, 118},
			SHA256:       [32]byte{179, 25, 234, 242, 132, 153, 97, 25, 128, 58, 41, 108, 187, 240, 247, 33, 90, 75, 41, 107, 253, 203, 214, 56, 60, 198, 238, 211, 119, 160, 221, 162},
			Homepage:     "http://example.com",
			Description:  "This is a test.\nIt does nothing.\n\nAbsolutely nothing.",
		},
	}

	want := `Package: appcast-test
Version: 1.1.0
Architecture: amd64
Maintainer: Test User <test@example.com>
Depends: bash
Recommends: git
Conflicts: appcast-test-new
Replaces: appcast-test-old
Provides: appcast-test-alt
Priority: optional
Section: utils
Filename: pool/main/a/appcast-test/appcast-test_1.1.0_amd64.deb
Size: 812
MD5sum: 0207b09e20c07bd5f2b05d99a7931919
SHA1: 0d34ab933af95a8574c7a9ce12da704e0f3f42eb
SHA256: f9b6cea0912e56fe22b7bef4a38897743a57ed9fe84e011bb375f033be6637f6
Homepage: http://example.com
Description: This is a test.
 It does nothing.
 .
 Absolutely nothing.

Package: appcast-test
Version: 1.0.0
Architecture: amd64
Maintainer: Test User <test@example.com>
Depends: bash
Recommends: git
Conflicts: appcast-test-new
Replaces: appcast-test-old
Provides: appcast-test-alt
Priority: optional
Section: utils
Filename: pool/main/a/appcast-test/appcast-test_1.0.0_amd64.deb
Size: 814
MD5sum: 1c15d1d66804cd073869e8bc223cf529
SHA1: 104c1ad2cfd828480e75947005dafc6257cf0576
SHA256: b319eaf284996119803a296cbbf0f7215a4b296bfdcbd6383cc6eed377a0dda2
Homepage: http://example.com
Description: This is a test.
 It does nothing.
 .
 Absolutely nothing.
`

	testMarshalUnmarshal(t, in, want)
}

func testMarshalUnmarshal[T any](t *testing.T, in T, want string) {
	t.Helper()

	t.Run("Marshal", func(t *testing.T) {
		b, err := deb.Marshal(in)
		if err != nil {
			t.Fatal(err)
		}

		if diff := cmp.Diff(want, string(b)); diff != "" {
			t.Fatal(diff)
		}
	})

	t.Run("Unmarshal", func(t *testing.T) {
		var v T
		if err := deb.Unmarshal([]byte(want), &v); err != nil {
			t.Fatal(err)
		}

		if diff := cmp.Diff(in, v); diff != "" {
			t.Fatal(diff)
		}
	})
}
