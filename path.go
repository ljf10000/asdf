package asdf

import (
	"os"
	"path/filepath"
	"strings"
)

func CurrentDirectory() string {
	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		return "."
	}

	return strings.Replace(dir, "\\", "/", -1)
}

func CurrentDirFile(file string) string {
	return filepath.Join(CurrentDirectory(), file)
}
