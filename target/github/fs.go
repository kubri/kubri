package github

import (
	"io/fs"
	"time"

	"github.com/google/go-github/github"
)

// githubDirEntry implements fs.DirEntry for GitHub repository content.
type githubDirEntry struct {
	content  *github.RepositoryContent
	modified func() time.Time
}

// Name returns the base name of the entry.
func (e *githubDirEntry) Name() string {
	return e.content.GetName()
}

// IsDir reports whether the entry is a directory.
func (e *githubDirEntry) IsDir() bool {
	return e.content.GetType() == "dir"
}

// Type returns a fs.FileMode; directories get fs.ModeDir.
func (e *githubDirEntry) Type() fs.FileMode {
	if e.IsDir() {
		return fs.ModeDir
	}
	return 0
}

// Info returns an fs.FileInfo for the entry.
func (e *githubDirEntry) Info() (fs.FileInfo, error) {
	return &githubFileInfo{e.content, e.modified}, nil
}

// githubFileInfo implements fs.FileInfo for GitHub repository content.
type githubFileInfo struct {
	content  *github.RepositoryContent
	modified func() time.Time
}

// Name returns the base name of the file.
func (fi *githubFileInfo) Name() string {
	return fi.content.GetName()
}

// Size returns the size of the file.
// For directories, Size is set to 0.
func (fi *githubFileInfo) Size() int64 {
	if fi.IsDir() {
		return 0
	}
	return int64(fi.content.GetSize())
}

// Mode returns file mode bits.
// Directories get fs.ModeDir with read/execute permissions,
// files are treated as read-only.
func (fi *githubFileInfo) Mode() fs.FileMode {
	if fi.IsDir() {
		return fs.ModeDir | 0o555
	}
	return 0o444
}

// ModTime returns the modification time.
func (fi *githubFileInfo) ModTime() time.Time {
	return fi.modified()
}

// IsDir reports whether the file is a directory.
func (fi *githubFileInfo) IsDir() bool {
	return fi.content.GetType() == "dir"
}

// Sys returns the underlying data source; here it is nil.
func (fi *githubFileInfo) Sys() any {
	return nil
}
