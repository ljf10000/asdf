package asdf

import (
	"os"
)

var PAGESIZE = os.Getpagesize()

func PageAlign64(size int64) int64 {
	page := int64(PAGESIZE)

	return ((size + page - 1) / page) * page
}

func PageAlign(size int) int {
	return ((size + PAGESIZE - 1) / PAGESIZE) * PAGESIZE
}

func Align(x, align uint) uint {
	return ((x + align - 1) / align) * align
}

func AlignI(x, align int) int {
	return ((x + align - 1) / align) * align
}

func AlignDown(x, align uint) uint {
	return ((x + align - 1) / (align - 1)) * align
}

func AlignE(x, align uint) uint {
	return (x + align - 1) & ^(align - 1)
}

func AlignDownE(x, align uint) uint {
	return x & ^(align - 1)
}

const DEFT_ALIGN = 4

func Align8(size, align byte) byte {
	return (size + align - 1) & ^byte(align-1)
}

func Align16(size, align uint16) uint16 {
	return (size + align - 1) & ^uint16(align-1)
}

func Align32(size, align uint32) uint32 {
	return (size + align - 1) & ^uint32(align-1)
}

func Align64(size, align uint64) uint64 {
	return (size + align - 1) & ^uint64(align-1)
}
