package asdf

import (
	"fmt"
)

type IdxZone struct {
	Begin uint32
	End   uint32
}

func (me *IdxZone) String() string {
	return fmt.Sprintf("begin(%d) end(%d)", me.Begin, me.End)
}

func (me *IdxZone) Zero() {
	me.Begin = 0
	me.End = 0
}

func (me *IdxZone) IsGood() bool {
	return me.Begin > 0 && me.End > 0
}

func (me *IdxZone) Count() uint32 {
	return me.End - me.Begin + 1
}

func (me *IdxZone) Diff() uint32 {
	return me.End - me.Begin
}

func (me *IdxZone) InZone(idx uint32) bool {
	return me.Begin <= idx && idx <= me.End
}

func (me *IdxZone) Match(v *IdxZone) bool {
	return me.InZone(v.Begin) || me.InZone(v.End)
}

func (me *IdxZone) Intersect(v *IdxZone) IdxZone {
	if me.InZone(v.Begin) {
		if me.InZone(v.End) {
			// |--------- me ---------|
			//     |----- v -----|
			return *v
		} else {
			// |--------- me ---------|
			//               |----- v -----|
			return IdxZone{
				Begin: v.Begin,
				End:   me.End,
			}
		}
	} else {
		if me.InZone(v.End) {
			//     |--------- me ---------|
			// |----- v -----|
			return IdxZone{
				Begin: me.Begin,
				End:   v.End,
			}

		} else {
			//                  |--------- me ---------|
			// |----- v -----|              or              |----- v -----|
			return IdxZone{}
		}
	}
}
