package asdf

import (
	"fmt"
	"testing"
)

func TestBinSprintf(t *testing.T) {
	bin := []byte("0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

	s := BinSprintf(bin)

	t.Logf("bin=%s", string(bin))
	t.Logf("s=\n%s", s)

	fmt.Printf("bin=%s\n", string(bin))
	fmt.Printf("s=\n%s\n", s)
}
