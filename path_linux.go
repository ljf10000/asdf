package asdf

import (
	"path/filepath"
)

func baseName(fullname string) string {
	return filepath.Base(fullname)
}
