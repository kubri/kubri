package test

import (
	"io/fs"
	"testing/fstest"
	"time"
)

// ReadFS reads the contents of a file system and returns a [fstest.MapFS].
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

// DirEntry is a static representation of a directory entry.
type DirEntry struct {
	Name  string
	IsDir bool
	Type  fs.FileMode
	Info  FileInfo
}

// FileInfo is a static representation of a file's information.
type FileInfo struct {
	Name    string
	Size    int64
	Mode    fs.FileMode
	ModTime time.Time
	IsDir   bool
	Sys     any
}

// ReadDirEntries takes a slice of [fs.DirEntry] and returns a slice of [DirEntry].
func ReadDirEntries(entries []fs.DirEntry) []DirEntry {
	files := make([]DirEntry, len(entries))

	for i, entry := range entries {
		info, err := entry.Info()
		if err != nil {
			panic(err)
		}
		files[i] = DirEntry{
			Name:  entry.Name(),
			IsDir: entry.IsDir(),
			Type:  entry.Type(),
			Info: FileInfo{
				Name:    info.Name(),
				Size:    info.Size(),
				Mode:    info.Mode(),
				ModTime: info.ModTime(),
				IsDir:   info.IsDir(),
				Sys:     info.Sys(),
			},
		}
	}

	return files
}
