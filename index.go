package asdf

import (
	"fmt"
)

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
	return me.Begin > 0 && me.End > 0
}

func (me *Idxzone) Count() int {
	return int(me.End-me.Begin) + 1
}

func (me *Idxzone) Diff() int {
	return int(me.End - me.Begin)
}

func (me *Idxzone) InZone(idx int) bool {
	return me.Begin <= uint32(idx) && uint32(idx) <= me.End
}

func (me *Idxzone) Match(v *Idxzone) bool {
	return me.InZone(int(v.Begin)) || me.InZone(int(v.End))
}

func (me *Idxzone) Intersect(v *Idxzone) Idxzone {
	if me.InZone(int(v.Begin)) {
		if me.InZone(int(v.End)) {
			// |--------- me ---------|
			//     |----- v -----|
			return *v
		} else {
			// |--------- me ---------|
			//               |----- v -----|
			return Idxzone{
				Begin: v.Begin,
				End:   me.End,
			}
		}
	} else {
		if me.InZone(int(v.End)) {
			//     |--------- me ---------|
			// |----- v -----|
			return Idxzone{
				Begin: me.Begin,
				End:   v.End,
			}

		} else {
			//                  |--------- me ---------|
			// |----- v -----|              or              |----- v -----|
			return Idxzone{}
		}
	}
}
