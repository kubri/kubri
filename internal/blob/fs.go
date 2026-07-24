package blob

import (
	"io/fs"
	"path/filepath"
	"time"

	"gocloud.dev/blob"
)

// blobDirEntry wraps *blob.ListObject and implements fs.DirEntry.
type blobDirEntry struct {
	obj *blob.ListObject
}

// Name returns the base name of the object key.
func (e *blobDirEntry) Name() string {
	return filepath.Base(e.obj.Key)
}

// IsDir reports whether the entry is a directory.
func (e *blobDirEntry) IsDir() bool {
	return e.obj.IsDir
}

// Type returns the file mode bits. If the entry is a directory,
// it returns fs.ModeDir; otherwise, it returns 0.
func (e *blobDirEntry) Type() fs.FileMode {
	if e.obj.IsDir {
		return fs.ModeDir
	}
	return 0
}

// Info returns a fs.FileInfo describing the file.
// Here we wrap our *blob.ListObject in a blobFileInfo.
func (e *blobDirEntry) Info() (fs.FileInfo, error) {
	return &blobFileInfo{e.obj}, nil
}

// blobFileInfo implements fs.FileInfo using data from *blob.ListObject.
type blobFileInfo struct {
	obj *blob.ListObject
}

// Name returns the base name of the object key.
func (fi *blobFileInfo) Name() string {
	return filepath.Base(fi.obj.Key)
}

// Size returns the size of the object.
func (fi *blobFileInfo) Size() int64 {
	return fi.obj.Size
}

// Mode returns the file mode bits. For directories, we add the fs.ModeDir flag.
func (fi *blobFileInfo) Mode() fs.FileMode {
	if fi.obj.IsDir {
		// Read and execute permissions for directories
		return fs.ModeDir | 0o555
	}
	// Read-only file permissions for regular files
	return 0o444
}

// ModTime returns the modification time of the object.
func (fi *blobFileInfo) ModTime() time.Time {
	return fi.obj.ModTime
}

// IsDir reports whether the entry is a directory.
func (fi *blobFileInfo) IsDir() bool {
	return fi.obj.IsDir
}

// Sys returns the underlying data source (nil in this case).
func (fi *blobFileInfo) Sys() any {
	return nil
}
