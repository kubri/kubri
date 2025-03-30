package target_test

import (
	"io"
	"io/fs"
	"testing"
	"testing/fstest"

	"github.com/google/go-cmp/cmp"

	"github.com/kubri/kubri/internal/test"
	"github.com/kubri/kubri/target"
	ftarget "github.com/kubri/kubri/target/file"
)

func TestCopyFS(t *testing.T) {
	dir := t.TempDir()
	tgt, _ := ftarget.New(ftarget.Config{Path: dir})

	fsys := fstest.MapFS{
		"file1.txt":     &fstest.MapFile{Data: []byte("content of file1")},
		"file2.txt":     &fstest.MapFile{Data: []byte("content of file2")},
		"dir/file3.txt": &fstest.MapFile{Data: []byte("content of file3")},
	}

	err := target.CopyFS(t.Context(), tgt, fsys)
	if err != nil {
		t.Fatalf("CopyFS failed: %v", err)
	}

	for path, file := range fsys {
		r, err := tgt.NewReader(t.Context(), path)
		if err != nil {
			t.Fatal(path, err)
		}
		defer r.Close()

		got, err := io.ReadAll(r)
		if err != nil {
			t.Fatal(path, err)
		}

		if diff := cmp.Diff(file.Data, got); diff != "" {
			t.Error(path, diff)
		}
	}
}

func TestCopyFSError(t *testing.T) {
	tests := []struct {
		name    string
		fsys    fs.FS
		wantErr error
	}{
		{
			name: "non-regular file",
			fsys: fstest.MapFS{
				"file1.txt": &fstest.MapFile{},
				"symlink":   &fstest.MapFile{Mode: fs.ModeSymlink},
			},
			wantErr: &fs.PathError{Op: "CopyFS", Path: "symlink", Err: fs.ErrInvalid},
		},
		{
			name: "error on open",
			fsys: fstest.MapFS{
				"..": &fstest.MapFile{},
			},
			wantErr: &fs.PathError{Op: "open", Path: "..", Err: fs.ErrNotExist},
		},
		{
			name: "error on read",
			fsys: &errFS{
				FS:      fstest.MapFS{"file1.txt": &fstest.MapFile{}},
				readErr: fs.ErrClosed,
			},
			wantErr: &fs.PathError{Op: "Copy", Path: "file1.txt", Err: fs.ErrClosed},
		},
	}

	opts := test.CompareErrorMessages()

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			dir := t.TempDir()
			tgt, _ := ftarget.New(ftarget.Config{Path: dir})

			err := target.CopyFS(t.Context(), tgt, test.fsys)
			if diff := cmp.Diff(test.wantErr, err, opts); diff != "" {
				t.Error(diff)
			}
		})
	}
}

type errFS struct {
	fs.FS
	readErr error
}

func (e *errFS) Open(name string) (fs.File, error) {
	if e.readErr != nil && name != "." {
		f, _ := e.FS.Open(name)
		return &errorFile{f, e.readErr}, nil
	}
	return e.FS.Open(name)
}

type errorFile struct {
	fs.File
	err error
}

func (e *errorFile) Read([]byte) (int, error) {
	return 0, e.err
}
