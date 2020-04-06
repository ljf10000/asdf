package asdf

import (
	"testing"
)

func TestBinSprintf(t *testing.T) {
	bin := []byte("0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

	s := BinSprintf(bin)

	t.Logf("bin=%s", string(bin))
	t.Logf("bin=\n%s", s)
}
