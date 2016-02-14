package emptyvfs

import (
	"fmt"
	"golang.org/x/tools/godoc/vfs"
	"os"
	"time"
)

// NewNameSpace returns a vfs.NameSpace initialized with an empty
// emulated directory mounted on the root mount point "/" so that
// directory traversal routines don't break if the user doesn't
// explicitly mount a FileSystem at "/". See the following issue:
// https://github.com/golang/go/issues/14190
func NewNameSpace() vfs.NameSpace {
	ns := vfs.NameSpace{}
	ns.Bind("/", newemptyVFS(), "/", vfs.BindReplace)
	return ns
}

// type emptyVFS emulates a vfs.FileSystem consisting of an empty directory
type emptyVFS struct {
	modtime time.Time
}

func newemptyVFS() *emptyVFS {
	return &emptyVFS{modtime: time.Now()}
}

// Open implements Opener. Since emptyVFS is an empty directory, all
// attempt to open a file returns errors.
func (e *emptyVFS) Open(path string) (vfs.ReadSeekCloser, error) {
	if path == "/" {
		return nil, fmt.Errorf("open: / is a directory")
	}
	return nil, os.ErrNotExist
}

// Stat returns FileInfo (implemented by emptyVFS itself) if the path
// is root "/" or error for any other path.
func (e *emptyVFS) Stat(path string) (os.FileInfo, error) {
	if path == "/" {
		return e, nil
	}
	return nil, os.ErrNotExist
}

// Lstat simply calls Stat
func (e *emptyVFS) Lstat(path string) (os.FileInfo, error) {
	return e.Stat(path)
}

// ReadDir returns an empty FileInfo slice if the path is root "/"
// else error for any other path.
func (e *emptyVFS) ReadDir(path string) ([]os.FileInfo, error) {
	if path == "/" {
		return []os.FileInfo{}, nil
	}
	return nil, os.ErrNotExist
}

func (e *emptyVFS) String() string {
	return "emptyVFS"
}

// These functions below implement os.FileInfo for the single
// empty emulated directory.

func (e *emptyVFS) Name() string {
	return "root"
}

func (e *emptyVFS) Size() int64 {
	return 0
}

func (e *emptyVFS) Mode() os.FileMode {
	return os.ModeDir | os.ModePerm
}

func (e *emptyVFS) ModTime() time.Time {
	return e.modtime
}

func (e *emptyVFS) IsDir() bool {
	return true
}

func (e *emptyVFS) Sys() interface{} {
	return nil
}
