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

func (me *SeqZone) IsGood() bool {
	return me.Begin > 0 && me.End > 0
}

func (me *SeqZone) Diff() uint64 {
	return me.End - me.Begin
}

func (me *SeqZone) Seq(idx uint64) uint64 {
	return me.Begin + idx
}

func (me *SeqZone) InZone(seq uint64) bool {
	return me.Begin <= seq && seq <= me.End
}

func (me *SeqZone) Match(v *SeqZone) bool {
	return me.InZone(v.Begin) || me.InZone(v.End)
}

func (me *SeqZone) Intersect(v *SeqZone) SeqZone {
	if me.InZone(v.Begin) {
		if me.InZone(v.End) {
			// |--------- me ---------|
			//     |----- v -----|
			return *v
		} else {
			// |--------- me ---------|
			//               |----- v -----|
			return SeqZone{
				Begin: v.Begin,
				End:   me.End,
			}
		}
	} else {
		if me.InZone(v.End) {
			//     |--------- me ---------|
			// |----- v -----|
			return SeqZone{
				Begin: me.Begin,
				End:   v.End,
			}

		} else {
			//                  |--------- me ---------|
			// |----- v -----|              or              |----- v -----|
			return SeqZone{}
		}
	}
}
