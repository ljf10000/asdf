package asdf

import (
	"testing"
)

func TestSize(t *testing.T) {
	ExecSizeCheck(sizeCheckers, true)
}
