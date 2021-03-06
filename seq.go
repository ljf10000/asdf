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

/******************************************************************************/

type Seq32 uint32

func (me Seq32) inRange(a, b Seq32) bool {
	return a <= me && me <= b
}

func (me Seq32) InRange(a, b Seq32) bool {
	if a < b {
		return me.inRange(a, b)
	} else {
		return me.inRange(b, a)
	}
}

func (me Seq32) InZone(z Seqzone32) bool {
	return me.inRange(z.Begin, z.End)
}

/******************************************************************************/

type Seqzone32 struct {
	Begin Seq32
	End   Seq32
}

func (me Seqzone32) String() string {
	return fmt.Sprintf("begin(%d) end(%d)", me.Begin, me.End)
}

func (me *Seqzone32) Zero() {
	me.Begin = 0
	me.End = 0
}

func (me Seqzone32) IsGood() bool {
	return me.Begin > 0 && me.End > 0
}

func (me Seqzone32) Count() Seq32 {
	return me.End - me.Begin + 1
}

func (me Seqzone32) Diff() Seq32 {
	return me.End - me.Begin
}

func (me Seqzone32) Include(z Seqzone32) bool {
	// |--------- me ---------|
	//      |----- z -----|
	return me.Begin <= z.Begin && me.End >= z.End
}

func (me Seqzone32) Compare(v Seqzone32) int {
	if me.End < v.Begin {
		// |--------- me ---------|
		//                            |----- v -----|
		return -1
	} else if me.Begin > v.End {
		//                  |--------- me ---------|
		// |----- v -----|
		return 1
	} else {
		//            |--------- me ---------|
		// |----- v -----|
		//                 |----- v -----|
		//                                 |----- v -----|
		return 0
	}
}

func (me Seqzone32) Match(v Seqzone32) bool {
	return 0 == me.Compare(v)
}

func (me Seqzone32) Intersect(v Seqzone32) (Seqzone32, bool) {
	if 0 != me.Compare(v) {
		return Seqzone32{}, false
	}

	// get max begin
	begin := me.Begin
	if me.Begin < v.Begin {
		begin = v.Begin
	}

	// get min end
	end := me.End
	if me.End > v.End {
		end = v.End
	}

	return Seqzone32{
		Begin: begin,
		End:   end,
	}, true
}

/******************************************************************************/
const SizeofSeqWraper = 8

type SeqWraper struct {
	High uint32
	Low  uint32
}

func (me *SeqWraper) Seq64() Seq64 {
	return Seq64(me.High)<<32 | Seq64(me.Low)
}

/******************************************************************************/

type Seq64 uint64

func (me Seq64) SeqWraper() SeqWraper {
	return SeqWraper{
		High: uint32(me >> 32),
		Low:  uint32(me & 0xffffffff),
	}
}

func (me Seq64) inRange(a, b Seq64) bool {
	return a <= me && me <= b
}

func (me Seq64) InRange(a, b Seq64) bool {
	if a < b {
		return me.inRange(a, b)
	} else {
		return me.inRange(b, a)
	}
}

func (me Seq64) InZone(z Seqzone) bool {
	return me.inRange(z.Begin, z.End)
}

/******************************************************************************/

type Seqzone struct {
	Begin Seq64
	End   Seq64
}

func (me Seqzone) String() string {
	return fmt.Sprintf("begin(%d) end(%d)", me.Begin, me.End)
}

func (me *Seqzone) Zero() {
	me.Begin = 0
	me.End = 0
}

func (me Seqzone) IsGood() bool {
	return me.Begin > 0 && me.End > 0
}

func (me Seqzone) Count() Seq64 {
	return me.End - me.Begin + 1
}

func (me Seqzone) Diff() Seq64 {
	return me.End - me.Begin
}

func (me Seqzone) Include(z Seqzone) bool {
	// |--------- me ---------|
	//      |----- z -----|
	return me.Begin <= z.Begin && me.End >= z.End
}

func (me Seqzone) Compare(v Seqzone) int {
	if me.End < v.Begin {
		// |--------- me ---------|
		//                            |----- v -----|
		return -1
	} else if me.Begin > v.End {
		//                  |--------- me ---------|
		// |----- v -----|
		return 1
	} else {
		//            |--------- me ---------|
		// |----- v -----|
		//                 |----- v -----|
		//                                 |----- v -----|
		return 0
	}
}

func (me Seqzone) Match(v Seqzone) bool {
	return 0 == me.Compare(v)
}

func (me Seqzone) Intersect(v Seqzone) (Seqzone, bool) {
	if 0 != me.Compare(v) {
		return Seqzone{}, false
	}

	// get max begin
	begin := me.Begin
	if me.Begin < v.Begin {
		begin = v.Begin
	}

	// get min end
	end := me.End
	if me.End > v.End {
		end = v.End
	}

	return Seqzone{
		Begin: begin,
		End:   end,
	}, true
}
