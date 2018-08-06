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

type Seqzone struct {
	Begin uint64
	End   uint64
}

func (me *Seqzone) String() string {
	return fmt.Sprintf("begin(%d) end(%d)", me.Begin, me.End)
}

func (me *Seqzone) Zero() {
	me.Begin = 0
	me.End = 0
}

func (me *Seqzone) IsGood() bool {
	return me.Begin > 0 && me.End > 0
}

func (me *Seqzone) Count() uint64 {
	return me.End - me.Begin + 1
}

func (me *Seqzone) Diff() uint64 {
	return me.End - me.Begin
}

func (me *Seqzone) InZone(seq uint64) bool {
	return me.Begin <= seq && seq <= me.End
}

func (me *Seqzone) Match(v Seqzone) bool {
	return me.InZone(v.Begin) || me.InZone(v.End)
}

func (me *Seqzone) Intersect(v Seqzone) Seqzone {
	if me.InZone(v.Begin) {
		if me.InZone(v.End) {
			// |--------- me ---------|
			//     |----- v -----|
			return v
		} else {
			// |--------- me ---------|
			//               |----- v -----|
			return Seqzone{
				Begin: v.Begin,
				End:   me.End,
			}
		}
	} else {
		if me.InZone(v.End) {
			//     |--------- me ---------|
			// |----- v -----|
			return Seqzone{
				Begin: me.Begin,
				End:   v.End,
			}

		} else {
			//                  |--------- me ---------|
			// |----- v -----|              or              |----- v -----|
			return Seqzone{}
		}
	}
}
