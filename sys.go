package asdf

import (
	"unsafe"
)

var _int int
var _int32 int32
var _int64 int64

func GetIntSize() int {
	if unsafe.Sizeof(_int) == unsafe.Sizeof(_int32) {
		return 4
	} else if unsafe.Sizeof(_int) == unsafe.Sizeof(_int64) {
		return 8
	} else {
		panic("unknow int size")
		return 0
	}
}
