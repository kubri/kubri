package target

import (
	"context"
	"io"
	"io/fs"
)

type Target interface {
	NewWriter(ctx context.Context, path string) (io.WriteCloser, error)
	NewReader(ctx context.Context, path string) (io.ReadCloser, error)
	Remove(ctx context.Context, path string) error
	Sub(dir string) Target
	URL(ctx context.Context, path string) (string, error)
}

// CopyFS copies the file system fsys to the target t.
func CopyFS(ctx context.Context, t Target, fsys fs.FS) error {
	return fs.WalkDir(fsys, ".", func(path string, d fs.DirEntry, err error) error {
		if err != nil || d.IsDir() {
			return err
		}

		// Cannot handle symlinks, devices, or other non-regular files.
		if !d.Type().IsRegular() {
			return &fs.PathError{Op: "CopyFS", Path: path, Err: fs.ErrInvalid}
		}

		r, err := fsys.Open(path)
		if err != nil {
			return err
		}
		defer r.Close()

		w, err := t.NewWriter(ctx, path)
		if err != nil {
			return err
		}

		if _, err := io.Copy(w, r); err != nil {
			w.Close()
			return &fs.PathError{Op: "Copy", Path: path, Err: err}
		}
		return w.Close()
	})
}
