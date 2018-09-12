package asdf

import (
	"os"
	"path/filepath"
)

func CurrentDirectory() string {
	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		return "."
	}

	return dir
}

func MakeDirectory(dir string) {
	// dirs := filepath.SplitList(dir)
}

func CurrentDirFile(file string) string {
	return filepath.Join(CurrentDirectory(), file)
}
