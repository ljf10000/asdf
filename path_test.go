package asdf

import (
	"path/filepath"
	"testing"
)

func TestPath(t *testing.T) {
	dirlist := []string{
		"/a/b/c",
		"/a/b/c/",
		"a/b/c",
		"a/b/c/",
	}

	for _, dir := range dirlist {
		dirs := filepath.SplitList(dir)

		t.Logf("%s==>%v\n", dir, dirs)

		abs, err := filepath.Abs(dir)
		if nil != err {
			t.Logf("abs(%s) error:%s\n", dir, err)
		}
		t.Logf("abs(%s)==>%s\n", dir, abs)
	}
}
