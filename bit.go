package asdf

type Flag8 uint8
type Flag16 uint16
type Flag32 uint32
type Flag64 uint64
type FlagU uint

/******************************************************************************/

func (me Flag8) Flag(bit int) Flag8 {
	return 1 << Flag8(bit)
}

func (me Flag16) Flag(bit int) Flag16 {
	return 1 << Flag16(bit)
}

func (me Flag32) Flag(bit int) Flag32 {
	return 1 << Flag32(bit)
}

func (me Flag64) Flag(bit int) Flag64 {
	return 1 << Flag64(bit)
}

func (me FlagU) Flag(bit int) FlagU {
	return 1 << FlagU(bit)
}

/******************************************************************************/
func (me Flag8) Set(flag Flag8) Flag8 {
	return me | flag
}

func (me Flag16) Set(flag Flag16) Flag16 {
	return me | flag
}

func (me Flag32) Set(flag Flag32) Flag32 {
	return me | flag
}

func (me Flag64) Set(flag Flag64) Flag64 {
	return me | flag
}

func (me FlagU) Set(flag FlagU) FlagU {
	return me | flag
}

/******************************************************************************/
func (me Flag8) SetBit(bit int) Flag8 {
	return me.Set(me.Flag(bit))
}

func (me Flag16) SetBit(bit int) Flag16 {
	return me.Set(me.Flag(bit))
}

func (me Flag32) SetBit(bit int) Flag32 {
	return me.Set(me.Flag(bit))
}

func (me Flag64) SetBit(bit int) Flag64 {
	return me.Set(me.Flag(bit))
}

func (me FlagU) SetBit(bit int) FlagU {
	return me.Set(me.Flag(bit))
}

/******************************************************************************/

func (me Flag8) Clear(flag Flag8) Flag8 {
	return me & ^flag
}

func (me Flag16) Clear(flag Flag16) Flag16 {
	return me & ^flag
}

func (me Flag32) Clear(flag Flag32) Flag32 {
	return me & ^flag
}

func (me Flag64) Clear(flag Flag64) Flag64 {
	return me & ^flag
}

func (me FlagU) Clear(flag FlagU) FlagU {
	return me & ^flag
}

/******************************************************************************/

func (me Flag8) ClearBit(bit int) Flag8 {
	return me.Clear(me.Flag(bit))
}

func (me Flag16) ClearBit(bit int) Flag16 {
	return me.Clear(me.Flag(bit))
}

func (me Flag32) ClearBit(bit int) Flag32 {
	return me.Clear(me.Flag(bit))
}

func (me Flag64) ClearBit(bit int) Flag64 {
	return me.Clear(me.Flag(bit))
}

func (me FlagU) ClearBit(bit int) FlagU {
	return me.Clear(me.Flag(bit))
}

/******************************************************************************/

func (me Flag8) Has(flag Flag8) bool {
	return flag == (flag & me)
}

func (me Flag16) Has(flag Flag16) bool {
	return flag == (flag & me)
}

func (me Flag32) Has(flag Flag32) bool {
	return flag == (flag & me)
}

func (me Flag64) Has(flag Flag64) bool {
	return flag == (flag & me)
}

func (me FlagU) Has(flag FlagU) bool {
	return flag == (flag & me)
}

/******************************************************************************/

func (me Flag8) HasBit(bit int) bool {
	return me.Has(me.Flag(bit))
}

func (me Flag16) HasBit(bit int) bool {
	return me.Has(me.Flag(bit))
}

func (me Flag32) HasBit(bit int) bool {
	return me.Has(me.Flag(bit))
}

func (me Flag64) HasBit(bit int) bool {
	return me.Has(me.Flag(bit))
}

func (me FlagU) HasBit(bit int) bool {
	return me.Has(me.Flag(bit))
}

//==============================================================================

func SetFlag8(x, flag uint8) uint8 {
	return x | flag
}
func SetFlag16(x, flag uint16) uint16 {
	return x | flag
}
func SetFlag32(x, flag uint32) uint32 {
	return x | flag
}
func SetFlag64(x, flag uint64) uint64 {
	return x | flag
}
func SetFlag(x, flag uint) uint {
	return x | flag
}

/******************************************************************************/
func ClrFlag8(x, flag uint8) uint8 {
	return x & ^flag
}
func ClrFlag16(x, flag uint16) uint16 {
	return x & ^flag
}
func ClrFlag32(x, flag uint32) uint32 {
	return x & ^flag
}
func ClrFlag64(x, flag uint64) uint64 {
	return x & ^flag
}
func ClrFlag(x, flag uint) uint {
	return x & ^flag
}

/******************************************************************************/
func HasFlag8(x, flag uint8) bool {
	return flag == (x & flag)
}
func HasFlag16(x, flag uint16) bool {
	return flag == (x & flag)
}
func HasFlag32(x, flag uint32) bool {
	return flag == (x & flag)
}
func HasFlag64(x, flag uint64) bool {
	return flag == (x & flag)
}
func HasFlag(x, flag uint) bool {
	return flag == (x & flag)
}

/******************************************************************************/
func SetBit8(x, bit uint8) uint8 {
	return SetFlag8(x, 1<<bit)
}
func SetBit16(x, bit uint16) uint16 {
	return SetFlag16(x, 1<<bit)
}
func SetBit32(x, bit uint32) uint32 {
	return SetFlag32(x, 1<<bit)
}
func SetBit64(x, bit uint64) uint64 {
	return SetFlag64(x, 1<<bit)
}
func SetBit(x, bit uint) uint {
	return SetFlag(x, 1<<bit)
}

/******************************************************************************/
func ClrBit8(x, bit uint8) uint8 {
	return ClrFlag8(x, 1<<bit)
}
func ClrBit16(x, bit uint16) uint16 {
	return ClrFlag16(x, 1<<bit)
}
func ClrBit32(x, bit uint32) uint32 {
	return ClrFlag32(x, 1<<bit)
}
func ClrBit64(x, bit uint64) uint64 {
	return ClrFlag64(x, 1<<bit)
}
func ClrBit(x, bit uint) uint {
	return ClrFlag(x, 1<<bit)
}

/******************************************************************************/
func HasBit8(x, bit uint8) bool {
	return HasFlag8(x, 1<<bit)
}
func HasBit16(x, bit uint16) bool {
	return HasFlag16(x, 1<<bit)
}
func HasBit32(x, bit uint32) bool {
	return HasFlag32(x, 1<<bit)
}
func HasBit64(x, bit uint64) bool {
	return HasFlag64(x, 1<<bit)
}
func HasBit(x, bit uint) bool {
	return HasFlag(x, 1<<bit)
}

/******************************************************************************/
type BitMap []uint32

const BitMapSlot = 32

func (me BitMap) isGoodIdx(idx uint32) bool {
	return int(idx) < len(me)
}

func (me BitMap) SetBit(bit uint32) {
	idx := bit / BitMapSlot

	if me.isGoodIdx(idx) {
		SetBit32(me[idx], bit%BitMapSlot)
	}
}

func (me BitMap) ClrBit(bit uint32) {
	idx := bit / BitMapSlot

	if me.isGoodIdx(idx) {
		ClrBit32(me[idx], bit%BitMapSlot)
	}
}

func (me BitMap) HasBit(bit uint32) bool {
	idx := bit / BitMapSlot

	if !me.isGoodIdx(idx) {
		return false
	}

	return HasBit32(me[idx], bit%BitMapSlot)
}

type BitsMapper struct {
	Type  string
	Bits  uint64
	Names map[int]string
}

func (me *BitsMapper) Name(flags uint64) string {
	name := make([]byte, 0, 1024)
	for i := uint64(0); i < me.Bits; i++ {
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
