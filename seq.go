package asdf

import (
	"fmt"
)

func Seq32Rand() uint32 {
	return RandSeed.Uint32() & 0xffff
}

func Seq32After(a, b uint32) bool {
	return (int32)(a-b) > 0
}

func Seq32Before(a, b uint32) bool {
	return (int32)(a-b) < 0
}

func Seq64Rand() uint64 {
	return uint64(RandSeed.Uint32())
}

func Seq64After(a, b uint64) bool {
	return (uint64)(a-b) > 0
}

func Seq64Before(a, b uint64) bool {
	return (uint64)(a-b) < 0
}

type SeqZone struct {
	Begin uint64
	End   uint64
}

func (me *SeqZone) String() string {
	return fmt.Sprintf("begin(%d) end(%d)", me.Begin, me.End)
}

func (me *SeqZone) Zero() {
	me.Begin = 0
	me.End = 0
}

func (me *SeqZone) Diff() uint64 {
	return me.End - me.Begin
}

func (me *SeqZone) Seq(idx uint64) uint64 {
	return me.Begin + idx
}
