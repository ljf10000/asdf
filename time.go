package asdf

import (
	"fmt"
	"time"
)

const TimeFormat = "2006-01-02 15:04:05"

/******************************************************************************/

func NowTime32() Time32 {
	return Time32(time.Now().Unix())
}

type Time32 uint32

func (me Time32) Unix() time.Time {
	return time.Unix(int64(me), 0)
}

func (me Time32) IsGood() bool {
	return me > 0
}

func (me Time32) Timespec() Timespec {
	return Timespec{
		Second: Time32(me),
		Nano:   0,
	}
}

func (me *Time32) Read(s string) error {
	tm, err := time.Parse(TimeFormat, s)
	if nil != err {
		return err
	}

	*me = Time32(tm.Unix())
	return nil
}

func (me Time32) inZone(a, b Time32) bool {
	return a <= me && me <= b
}

func (me Time32) InZone(a, b Time32) bool {
	if a < b {
		return me.inZone(a, b)
	} else {
		return me.inZone(b, a)
	}
}

/******************************************************************************/
type Timezone32 struct {
	Begin Time32
	End   Time32
}

func (me *Timezone32) String() string {
	return fmt.Sprintf("begin(%s) end(%s)", me.Begin, me.End)
}

func (me *Timezone32) IsGood() bool {
	return me.Begin.IsGood() && me.End.IsGood()
}

func (me *Timezone32) Zero() {
	me.Begin = 0
	me.End = 0
}

func (me *Timezone32) InZone(v Time32) bool {
	return v.inZone(me.Begin, me.End)
}

/******************************************************************************/

type Timens = Time32
type Time64 uint64

func MakeTime64(second Time32, nano Timens) Time64 {
	return Time64(second)*1e9 + Time64(nano)
}

func (me Time64) Eq(v Time64) bool {
	return me == v
}

func (me Time64) Le(v Time64) bool {
	return me < v
}

func (me Time64) Ge(v Time64) bool {
	return me > v
}

func (me Time64) Compare(v Time64) (int, Time64 /*diff*/) {
	switch {
	case me > v:
		return 1, me - v
	case me == v:
		return 0, 0
	default: // case me < v:
		return -1, v - me
	}
}

func (me Time64) Split() (Time32, Timens) {
	return Time32(me / 1e9), Timens(me % 1e9)
}

func (me Time64) Timespec() Timespec {
	return Timespec{
		Second: Time32(me / 1e9),
		Nano:   Timens(me % 1e9),
	}
}

/******************************************************************************/

type Timespec struct {
	Second Time32
	Nano   Timens
}

func MakeTimespec(second Time32, nano Timens) Timespec {
	return Timespec{
		Second: second,
		Nano:   nano,
	}
}

func (me *Timespec) IsGood() bool {
	return me.Second > 0 || me.Nano > 0
}

func (me *Timespec) Zero() {
	me.Second = 0
	me.Nano = 0
}

func (me *Timespec) Time64() Time64 {
	return MakeTime64(me.Second, me.Nano)
}

func (me *Timespec) String() string {
	t := time.Unix(int64(me.Second), int64(me.Nano))

	return t.String()
}

func (me *Timespec) Load(t Time64) {
	me.Second, me.Nano = t.Split()
}

func (me *Timespec) Compare(v Timespec) (int, Timespec /*diff*/) {
	a := me.Time64()
	b := me.Time64()

	cmp, diff := a.Compare(b)

	return cmp, diff.Timespec()
}

func (me *Timespec) inZone(a, b Timespec) bool {
	if cmp, _ := me.Compare(a); cmp < 0 { // me < a
		return false
	}

	if cmp, _ := me.Compare(b); cmp > 0 { // me > b
		return false
	}

	return true
}

func (me *Timespec) InZone(a, b Timespec) bool {
	cmp, _ := a.Compare(b)
	if cmp < 0 {
		// a < b
		return me.inZone(a, b)
	} else {
		// a >= b
		return me.inZone(b, a)
	}
}

func (me *Timespec) Eq(v Timespec) bool {
	return *me == v
}

func (me *Timespec) Le(v Timespec) bool {
	cmp, _ := me.Compare(v)

	return cmp < 0
}

func (me *Timespec) Ge(v Timespec) bool {
	cmp, _ := me.Compare(v)

	return cmp > 0
}

func (me *Timespec) Add(v Timespec) Timespec {
	t := me.Time64() + v.Time64()

	return t.Timespec()
}

/******************************************************************************/

type Timezone struct {
	Begin Timespec
	End   Timespec
}

func (me *Timezone) String() string {
	return fmt.Sprintf("begin(%s) end(%s)", me.Begin, me.End)
}

func (me *Timezone) IsGood() bool {
	return me.Begin.IsGood() && me.End.IsGood()
}

func (me *Timezone) Timezone32() Timezone32 {
	return Timezone32{
		Begin: me.Begin.Second,
		End:   me.End.Second,
	}
}

func (me *Timezone) Zero() {
	me.Begin.Zero()
	me.End.Zero()
}

func (me *Timezone) InZone32(v Time32) bool {
	return v.inZone(me.Begin.Second, me.End.Second)
}

func (me *Timezone) InZone(v Timespec) bool {
	return v.inZone(me.Begin, me.End)
}

func (me *Timezone) Match32(v Timezone32) bool {
	return me.InZone32(v.Begin) || me.InZone32(v.End)
}

func (me *Timezone) Match(v *Timezone) bool {
	return me.InZone(v.Begin) || me.InZone(v.End)
}

func (me *Timezone) Intersect32(v Timezone32) Timezone32 {
	if me.InZone32(v.Begin) {
		if me.InZone32(v.End) {
			// |--------- me ---------|
			//     |----- v -----|
			return v
		} else {
			// |--------- me ---------|
			//               |----- v -----|
			return Timezone32{
				Begin: v.Begin,
				End:   me.End.Second,
			}
		}
	} else {
		if me.InZone32(v.End) {
			//     |--------- me ---------|
			// |----- v -----|
			return Timezone32{
				Begin: me.Begin.Second,
				End:   v.End,
			}

		} else {
			//                  |--------- me ---------|
			// |----- v -----|              or              |----- v -----|
			return Timezone32{}
		}
	}
}

func (me *Timezone) Intersect(v *Timezone) Timezone {
	if me.InZone(v.Begin) {
		if me.InZone(v.End) {
			// |--------- me ---------|
			//     |----- v -----|
			return *v
		} else {
			// |--------- me ---------|
			//               |----- v -----|
			return Timezone{
				Begin: v.Begin,
				End:   me.End,
			}
		}
	} else {
		if me.InZone(v.End) {
			//     |--------- me ---------|
			// |----- v -----|
			return Timezone{
				Begin: me.Begin,
				End:   v.End,
			}

		} else {
			//                  |--------- me ---------|
			// |----- v -----|              or              |----- v -----|
			return Timezone{}
		}
	}
}
