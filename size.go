package asdf

import (
	"fmt"
)

type SizeChecker struct {
	Name  string
	Size  uintptr // unsafe.Sizeof
	SizeC int     // const size
}

func NewSizeChecker(name string, size uintptr, csize int) *SizeChecker {
	return &SizeChecker{
		Name:  name,
		Size:  size,
		SizeC: csize,
	}
}

func (me *SizeChecker) Exec(show bool) {
	if me.SizeC != int(me.Size) {
		Panic("%s self check failed size(%d) != calc-size(%d)", me.Name, me.SizeC, me.Size)
	}

	if show {
		fmt.Printf("sizeof(%s)=%d\n", me.Name, me.Size)
	}
}

// do it when init
func ExecSizeCheck(scs []*SizeChecker, show bool) {
	for i := 0; i < len(scs); i++ {
		scs[i].Exec(show)
	}
}
