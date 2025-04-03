package test

import (
	"io/fs"
	"testing/fstest"
)

// ReadFS reads the contents of a file system and returns a map of file paths to
// their contents.
func ReadFS(fsys fs.FS) fstest.MapFS {
	files := fstest.MapFS{}

	err := fs.WalkDir(fsys, ".", func(p string, d fs.DirEntry, err error) error {
		if err != nil || d.IsDir() {
			return err
		}
		content, err := fs.ReadFile(fsys, p)
		if err != nil {
			return err
		}
		fi, err := d.Info()
		if err != nil {
			return err
		}
		files[p] = &fstest.MapFile{
			Data:    content,
			Mode:    fi.Mode(),
			ModTime: fi.ModTime(),
			Sys:     fi.Sys(),
		}
		return nil
	})
	if err != nil {
		panic(err)
	}

	return files
}
