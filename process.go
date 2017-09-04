package asdf

import (
	"os"
)

var StdAttr = &os.ProcAttr{Files: []*os.File{os.Stdin, os.Stdout, os.Stderr}}
