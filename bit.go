package asdf

type Flag32 uint32
type Flag64 uint64

const (
	FlagOne32 Flag32 = 1
	FlagOne64 Flag64 = 1
)

func (me Flag32) Flag(bit int) Flag32 {
	return Flag32(1) << Flag32(bit)
}

func (me Flag64) Flag(bit int) Flag64 {
	return Flag64(1) << Flag64(bit)
}

func (me Flag32) Set(flag Flag32) Flag32 {
	return me | flag
}

func (me Flag64) Set(flag Flag64) Flag64 {
	return me | flag
}

func (me Flag32) SetBit(bit int) Flag32 {
	return me.Set(me.Flag(bit))
}

func (me Flag64) SetBit(bit int) Flag64 {
	return me.Set(me.Flag(bit))
}

func (me Flag32) Clear(flag Flag32) Flag32 {
	return me & ^flag
}

func (me Flag64) Clear(flag Flag64) Flag64 {
	return me & ^flag
}

func (me Flag32) ClearBit(bit int) Flag32 {
	return me.Clear(me.Flag(bit))
}

func (me Flag64) ClearBit(bit int) Flag64 {
	return me.Clear(me.Flag(bit))
}

func (me Flag32) Has(flag Flag32) bool {
	return flag == (flag & me)
}

func (me Flag64) Has(flag Flag64) bool {
	return flag == (flag & me)
}

func (me Flag64) HasBit(bit int) bool {
	return me.Has(me.Flag(bit))
}

func (me Flag32) HasBit(bit int) bool {
	return me.Has(me.Flag(bit))
}

//==============================================================================

func SetFlag(x, flag uint32) uint32 {
	return x | flag
}

func ClrFlag(x, flag uint32) uint32 {
	return x & ^flag
}

func HasFlag(x, flag uint32) bool {
	return flag == (x & flag)
}

func SetBit(x, bit uint32) uint32 {
	return SetFlag(x, 1<<bit)
}

func ClrBit(x, bit uint32) uint32 {
	return ClrFlag(x, 1<<bit)
}

func HasBit(x, bit uint32) bool {
	return HasFlag(x, 1<<bit)
}

type BitMap []uint32

const BitMapSlot = 32

func (me BitMap) isGoodIdx(idx uint32) bool {
	return int(idx) < len(me)
}

func (me BitMap) SetBit(bit uint32) {
	idx := bit / BitMapSlot

	if me.isGoodIdx(idx) {
		SetBit(me[idx], bit%BitMapSlot)
	}
}

func (me BitMap) ClrBit(bit uint32) {
	idx := bit / BitMapSlot

	if me.isGoodIdx(idx) {
		ClrBit(me[idx], bit%BitMapSlot)
	}
}

func (me BitMap) HasBit(bit uint32) bool {
	idx := bit / BitMapSlot

	if !me.isGoodIdx(idx) {
		return false
	}

	return HasBit(me[idx], bit%BitMapSlot)
}

type BitsMapper struct {
	Type  string
	Names map[int]string
}

func (me *BitsMapper) name(flags, bits uint64) string {
	name := make([]byte, 0, 1024)
	for i := uint64(0); i < bits; i++ {
		flag := uint64(1) << i

		if flag == (flag & flags) {
			v, ok := me.Names[int(flag)]
			if ok {
				name = append(name, []byte(v)...)
				name = append(name, '|')
			}
		}
	}

	Len := len(name)
	if Len > 0 {
		return string(name[:Len-1])
	} else {
		return Unknow
	}
}

func (me *BitsMapper) Name8(flags int) string {
	return me.name(uint64(flags), 8*SizeofByte)
}

func (me *BitsMapper) Name16(flags int) string {
	return me.name(uint64(flags), 8*SizeofInt16)
}

func (me *BitsMapper) Name32(flags int) string {
	return me.name(uint64(flags), 8*SizeofInt32)
}

func (me *BitsMapper) Name64(flags int) string {
	return me.name(uint64(flags), 8*SizeofInt64)
}
