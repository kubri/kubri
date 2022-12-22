package local_test

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/abemedia/appcast/source"
	"github.com/abemedia/appcast/source/local"
	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
)

func TestLocal(t *testing.T) {
	dir := t.TempDir()

	tests := []struct {
		name string
		path string
		want []*source.Release
	}{
		{
			name: "Directory",
			path: filepath.Join(dir, "dir"),
			want: []*source.Release{
				{
					Name:    "v0.0.0",
					Version: "v0.0.0",
					Date:    time.Now(),
					Assets: []*source.Asset{
						{
							Name: "test.dmg",
							URL:  "file://" + filepath.Join(dir, "dir", "test.dmg"),
							Size: 5,
						},
						{
							Name: "test_32-bit.msi",
							URL:  "file://" + filepath.Join(dir, "dir", "test_32-bit.msi"),
							Size: 5,
						},
						{
							Name: "test_64-bit.msi",
							URL:  "file://" + filepath.Join(dir, "dir", "test_64-bit.msi"),
							Size: 5,
						},
					},
				},
			},
		},
		{
			name: "Glob",
			path: filepath.Join(dir, "glob", "*", "*"),
			want: []*source.Release{
				{
					Name:    "v0.0.0",
					Version: "v0.0.0",
					Date:    time.Now(),
					Assets: []*source.Asset{
						{
							Name: "test.dmg",
							URL:  "file://" + filepath.Join(dir, "glob", "darwin", "test.dmg"),
							Size: 5,
						},
						{
							Name: "test_32-bit.msi",
							URL:  "file://" + filepath.Join(dir, "glob", "win32", "test_32-bit.msi"),
							Size: 5,
						},
						{
							Name: "test_64-bit.msi",
							URL:  "file://" + filepath.Join(dir, "glob", "win64", "test_64-bit.msi"),
							Size: 5,
						},
					},
				},
			},
		},
		{
			name: "File",
			path: filepath.Join(dir, "file", "test.dmg"),
			want: []*source.Release{
				{
					Name:    "v0.0.0",
					Version: "v0.0.0",
					Date:    time.Now(),
					Assets: []*source.Asset{
						{
							Name: "test.dmg",
							URL:  "file://" + filepath.Join(dir, "file", "test.dmg"),
							Size: 5,
						},
					},
				},
			},
		},
	}

	data := []byte("test\n")

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			for _, release := range test.want {
				for _, asset := range release.Assets {
					path := strings.TrimPrefix(asset.URL, "file://")
					os.MkdirAll(filepath.Dir(path), os.ModePerm)
					os.WriteFile(path, data, os.ModePerm)
				}
			}

			s, err := local.New(source.Config{Repo: test.path})
			if err != nil {
				t.Fatal(err)
			}

			opt := cmpopts.EquateApproxTime(100 * time.Millisecond)

			t.Run("ListReleases", func(t *testing.T) {
				got, err := s.ListReleases(nil)
				if err != nil {
					t.Fatal(err)
				}

				if diff := cmp.Diff(test.want, got, opt); diff != "" {
					t.Error(diff)
				}
			})

			t.Run("GetRelease", func(t *testing.T) {
				got, err := s.GetRelease(test.want[0].Version)
				if err != nil {
					t.Fatal(err)
				}

				if diff := cmp.Diff(test.want[0], got, opt); diff != "" {
					t.Error(diff)
				}
			})

			t.Run("UploadAsset", func(t *testing.T) {
				err := s.UploadAsset(test.want[0].Version, "test.txt", data)
				if err != nil {
					t.Fatal(err)
				}

				path, _, _ := strings.Cut(test.path, "*")
				if fi, _ := os.Stat(path); !fi.IsDir() {
					path = filepath.Dir(path)
				}

				b, err := os.ReadFile(filepath.Join(path, "test.txt"))
				if err != nil {
					t.Fatal(err)
				}

				if !bytes.Equal(data, b) {
					t.Error("should be equal")
				}
			})

			t.Run("DownloadAsset", func(t *testing.T) {
				b, err := s.DownloadAsset(test.want[0].Version, test.want[0].Assets[0].Name)
				if err != nil {
					t.Fatal(err)
				}

				if !bytes.Equal(data, b) {
					t.Error("should be equal")
				}
			})
		})
	}
}
