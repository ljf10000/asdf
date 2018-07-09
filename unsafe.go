package asdf

import (
	"fmt"
	"reflect"
	"unsafe"
)

func MakeTimespec(second, nano uint32) uint64 {
	return (uint64(second) << 32) | uint64(nano)
}

func SplitTimespec(timespec uint64) (uint32, uint32) {
	return uint32(timespec >> 32), uint32(timespec & 0xffffffff)
}

func MakeSlice(Data uintptr, Len, Cap int) []byte {
	var s = reflect.SliceHeader{
		Data: Data,
		Len:  Len,
		Cap:  Cap,
	}

	return *(*[]byte)(unsafe.Pointer(&s))
}

func SliceAddress(buf []byte) uintptr {
	return ((*reflect.SliceHeader)(unsafe.Pointer(&buf))).Data
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
		panic(fmt.Sprintf("%s self check failed", me.Name))
	}
}

func ExecSizeCheck(scs []*SizeChecker, show bool) {
	for i := 0; i < len(scs); i++ {
		scs[i].Exec(show)
	}
}
