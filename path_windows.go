package asdf

import (
	"path/filepath"
	"strings"
)

func baseName(fullname string) string {
	name := filepath.Base(fullname)

	split := strings.Split(name, FileNameSplit)

	return split[0]
}
