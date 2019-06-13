package asdf

import (
	"unsafe"
)

const (
	FLAT_STRING_ALIGN      = 1 << FLAT_STRING_ALIGN_BITS
	FLAT_STRING_ALIGN_BITS = 3
	FLAT_STRING_ALIGN_MASK = FLAT_STRING_ALIGN - 1
)

func flatStringAlign(v int) int {
	return (v + FLAT_STRING_ALIGN - 1) & ^FLAT_STRING_ALIGN_MASK
}

var (
	zFlatString  = FlatString{}
	scFlatString = NewSizeChecker("FlatString", unsafe.Sizeof(zFlatString), SizeofFlatString)
)

const SizeofFlatString = SizeofInt32

// flat string is a buffer
// 	1. size, 4 byte
//  2. content, n byte
type FlatString struct {
	size uint32
}

func (me *FlatString) body() uintptr {
	return uintptr(unsafe.Pointer(me)) + SizeofFlatString
}

func (me *FlatString) makeSlice() []byte {
	size := int(me.size)

	return MakeSliceEx(me.body(), size, size)
}

func (me *FlatString) makeString() string {
	return MakeStringEx(me.body(), int(me.size))
}

func (me *FlatString) Write(s string) {
	me.size = uint32(len(s))

	copy(me.makeSlice(), s)
}

func (me *FlatString) Len() int {
	size := int(me.size)

	return flatStringAlign(size) + SizeofFlatString
}

func (me *FlatString) Bin() []byte {
	if nil != me {
		return me.makeSlice()
	} else {
		return nil
	}
}

func (me *FlatString) String() string {
	if nil != me {
		return me.makeString()
	} else {
		return Empty
	}
}
