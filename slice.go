package asdf

import (
	"reflect"
	"unsafe"
)

type Slice []byte

func (me *Slice) Header() *reflect.SliceHeader {
	return (*reflect.SliceHeader)(unsafe.Pointer(me))
}

func (me Slice) IsValue(Value byte) bool {
	for i := 0; i < len(me); i++ {
		if Value != me[i] {
			return false
		}
	}

	return true
}

func (me Slice) IsZero() bool {
	return me.IsValue(0)
}

func (me Slice) IsFull() bool {
	return me.IsValue(255)
}

func (me Slice) Slice() []byte {
	return me
}

func (me Slice) Eq(it interface{}) bool {
	v, ok := it.(ISlice)
	if !ok {
		return false
	}
	b := v.Slice()

	if len(me) != len(b) {
		return false
	}

	for i := 0; i < len(me); i++ {
		if me[i] != b[i] {
			return false
		}
	}

	return true
}

func ObjToSlice(obj unsafe.Pointer, size int) []byte {
	if nil != obj {
		bin := Slice{}
		hdr := bin.Header()
		hdr.Data = uintptr(obj)
		hdr.Len = size
		hdr.Cap = size

		return bin
	} else {
		return nil
	}
}

func MemberToSlice(obj unsafe.Pointer, offset, size uintptr) []byte {
	member := unsafe.Pointer(uintptr(obj) + offset)

	return ObjToSlice(member, int(size))
}

func MakeSlice(Data uintptr, Len, Cap int) []byte {
	s := Slice{}
	h := s.Header()

	h.Data = Data
	h.Len = Len
	h.Cap = Cap

	return s
}

func StructSlice(data uintptr, size int) []byte {
	return MakeSlice(data, size, size)
}

func SliceAddress(buf []byte) uintptr {
	return ((*reflect.SliceHeader)(unsafe.Pointer(&buf))).Data
}

func SlicePointer(buf []byte) unsafe.Pointer {
	return unsafe.Pointer(((*reflect.SliceHeader)(unsafe.Pointer(&buf))).Data)
}
