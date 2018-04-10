package asdf

import (
	"testing"
)

func TestBkdr(t *testing.T) {
	s := "001xugoU0YUiYX1qkRnU0pOmoU0xugoP"

	bkdr := GenBkdr([]byte(s))

	t.Log("s=", s, ", bkdr=", bkdr)
}
