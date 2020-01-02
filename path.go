package asdf

import (
	"os"
	"path/filepath"
)

func CurrentDirectory() string {
	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		return FileNameSplit
	}

	return dir
}

func CurrentDirFile(file string) string {
	return filepath.Join(CurrentDirectory(), file)
}

func WorkDirectory() string {
	return CurrentDirectory()
}
