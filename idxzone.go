package asdf

import (
	"fmt"
)

// [begin, end]
type Idxzone struct {
	Begin uint32
	End   uint32
}

func (me *Idxzone) String() string {
	return fmt.Sprintf("begin(%d) end(%d)", me.Begin, me.End)
}

func (me *Idxzone) Zero() {
	me.Begin = 0
	me.End = 0
}

func (me *Idxzone) IsGood() bool {
	return me.Begin >= 0 && me.End > me.Begin
}

func (me *Idxzone) Count() int {
	return int(me.End-me.Begin) + 1
}

func (me *Idxzone) Diff() int {
	return int(me.End - me.Begin)
}

func (me *Idxzone) Include(z Idxzone) bool {
	// |--------- me ---------|
	//      |----- z -----|
	return me.Begin <= z.Begin && me.End >= z.End
}

func (me *Idxzone) InZone(idx int) bool {
	return me.Begin <= uint32(idx) && uint32(idx) <= me.End
}

func (me *Idxzone) Compare(v Idxzone) int {
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
		//       |--------------- v ----------------|
		return 0
	}
}

func (me *Idxzone) Match(v Idxzone) bool {
	return 0 == me.Compare(v)
}

func (me *Idxzone) Intersect(v Idxzone) Idxzone {
	if 0 != me.Compare(v) {
		return Idxzone{}
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

	return Idxzone{
		Begin: begin,
		End:   end,
	}
}
