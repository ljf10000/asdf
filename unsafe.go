package asdf

import (
	"fmt"
)

func MakeTimespec(second, nano uint32) uint64 {
	return (uint64(second) << 32) | uint64(nano)
}

func SplitTimespec(timespec uint64) (uint32, uint32) {
	return uint32(timespec >> 32), uint32(timespec & 0xffffffff)
}

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
		s := fmt.Sprintf("%s self check failed xsize(%d)!=size(%d)", me.Name, me.SizeX, me.Size)

		panic(s)
	}
}

func ExecSizeCheck(scs []*SizeChecker, show bool) {
	for i := 0; i < len(scs); i++ {
		scs[i].Exec(show)
	}
}
