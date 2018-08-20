package asdf

import (
	"testing"
)

func TestSliceCopy(t *testing.T) {
	s := "whats the fuck"
	var d []byte

	// when dst is nil
	// 	copy not panic, but NOT copy anything
	copy(d, []byte(s))

	if nil == d {
		t.Logf("d is nil")
	} else {
		t.Logf("d=%s", string(d))
	}

}
