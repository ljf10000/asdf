package asdf

import (
	"fmt"
)

type SizeChecker struct {
	Name  string
	Size  uintptr
	SizeX int
}

func NewSizeChecker(name string, size uintptr, sizex int) *SizeChecker {
	return &SizeChecker{
		Name:  name,
		Size:  size,
		SizeX: sizex,
	}
}

func (me *SizeChecker) Exec(show bool) {
	if show {
		fmt.Printf("sizeof(%s)=%d\n", me.Name, me.Size)
	}

	if me.SizeX != int(me.Size) {
		s := fmt.Sprintf("%s self check failed size(%d) != calc-size(%d)", me.Name, me.SizeX, me.Size)

		panic(s)
	}
}

func ExecSizeCheck(scs []*SizeChecker, show bool) {
	for i := 0; i < len(scs); i++ {
		scs[i].Exec(show)
	}
}
