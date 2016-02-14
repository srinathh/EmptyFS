package emptyvfs_test

import (
	"github.com/srinathh/emptyvfs"
	"golang.org/x/tools/godoc/vfs"
	"golang.org/x/tools/godoc/vfs/mapfs"
	"testing"
)

func TestNewNameSpace(t *testing.T) {

	// We will mount this filesystem under /fs1
	mount := mapfs.New(map[string]string{"fs1file": "abcdefgh"})

	// Existing process. This should give error on Stat("/")
	t1 := vfs.NameSpace{}
	t1.Bind("/fs1", mount, "/", vfs.BindReplace)

	// using NewNameSpace. This should work fine.
	t2 := emptyvfs.NewNameSpace()
	t2.Bind("/fs1", mount, "/", vfs.BindReplace)

	testcases := map[string][]bool{
		"/":            []bool{false, true},
		"/fs1":         []bool{true, true},
		"/fs1/fs1file": []bool{true, true},
	}

	fss := []vfs.FileSystem{t1, t2}

	for j, fs := range fss {
		for k, v := range testcases {
			_, err := fs.Stat(k)
			result := err == nil
			if result != v[j] {
				t.Errorf("fs: %d, testcase: %s, want: %v, got: %v, err: %s", j, k, v[j], result, err)
			}
		}
	}
}
