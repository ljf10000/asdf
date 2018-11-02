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

func (me Time32) String() string {
	return me.Unix().Format(TimeFormat)
}

func (me Time32) Unix() time.Time {
	return time.Unix(int64(me), 0)
}

func (me Time32) ToUnix() time.Time {
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

func (me Time32) inRange(a, b Time32) bool {
	return a <= me && me <= b
}

func (me Time32) InRange(a, b Time32) bool {
	if a < b {
		return me.inRange(a, b)
	} else {
		return me.inRange(b, a)
	}
}

func (me Time32) InZone(z Timezone32) bool {
	return me.inRange(z.Begin, z.End)
}

/******************************************************************************/

type Timens = Time32
type Timeus = Time32
type Timems = Time32

type Time64 uint64

func NowTime64() Time64 {
	return Time64(time.Now().UnixNano())
}

func MakeTime64(second Time32, nano Timens) Time64 {
	return Time64(second)*1e9 + Time64(nano)
}

func (me Time64) String() string {
	return me.Unix().Format(TimeFormat)
}

func (me Time64) Diff(tm Time64) int {
	if me > tm {
		return int(me - tm)
	} else {
		return int(tm - me)
	}
}

func (me Time64) Unix() time.Time {
	s, n := me.Split()

	return time.Unix(int64(s), int64(n))
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

func (me *Time64) Read(s string) error {
	tm, err := time.Parse(TimeFormat, s)
	if nil != err {
		return err
	}

	*me = Time64(tm.UnixNano())

	return nil
}

func (me Time64) Timeval() Timeval {
	return Timeval{
		Second: Time32(me / 1e9),
		Micro:  1000 * Timens(me%1e9),
	}
}

/******************************************************************************/
const SizeofTimespec = 8

type Timespec struct {
	Second Time32 `json:"second"`
	Nano   Timens `json:"nano"`
}

func MakeTimespec(second Time32, nano Timens) Timespec {
	return Timespec{
		Second: second,
		Nano:   nano,
	}
}

func (me Timespec) Diff(tm Timespec) int {
	return me.Time64().Diff(tm.Time64())
}

func (me Timespec) IsGood() bool {
	return me.Second > 0
}

func (me Timespec) Zero() {
	me.Second = 0
	me.Nano = 0
}

func (me Timespec) Up() Timespec {
	return Timespec{
		Second: me.Second + 1,
		Nano:   0,
	}
}

func (me Timespec) Down() Timespec {
	return Timespec{
		Second: me.Second,
		Nano:   0,
	}
}

func (me Timespec) Time64() Time64 {
	return MakeTime64(me.Second, me.Nano)
}

func (me Timespec) Timeval() Timeval {
	return Timeval{
		Second: me.Second,
		Micro:  me.Nano / 1000,
	}
}

func (me Timespec) String() string {
	return time.Unix(int64(me.Second), int64(me.Nano)).Format(TimeFormat)
}

func (me Timespec) Load(t Time64) {
	me.Second, me.Nano = t.Split()
}

func (me Timespec) Compare(v Timespec) (int, Timespec /*diff*/) {
	a := me.Time64()
	b := v.Time64()

	cmp, diff := a.Compare(b)

	return cmp, diff.Timespec()
}

func (me Timespec) inRange(a, b Timespec) bool {
	if cmp, _ := me.Compare(a); cmp < 0 { // me < a
		return false
	}

	if cmp, _ := me.Compare(b); cmp > 0 { // me > b
		return false
	}

	return true
}

func (me Timespec) InRange(a, b Timespec) bool {
	cmp, _ := a.Compare(b)
	if cmp < 0 {
		// a < b
		return me.inRange(a, b)
	} else {
		// a >= b
		return me.inRange(b, a)
	}
}

func (me Timespec) InZone(z Timezone) bool {
	return me.inRange(z.Begin, z.End)
}

func (me Timespec) Add(v Timespec) Timespec {
	t := me.Time64() + v.Time64()

	return t.Timespec()
}

/******************************************************************************/

type Timeval struct {
	Second Time32 `json:"second"`
	Micro  Timeus `json:"micro"`
}

func MakeTimeval(second Time32, micro Timeus) Timeval {
	return Timeval{
		Second: second,
		Micro:  micro,
	}
}

func (me Timeval) IsGood() bool {
	return me.Second > 0 || me.Micro > 0
}

func (me *Timeval) Zero() {
	me.Second = 0
	me.Micro = 0
}

func (me Timeval) Up() Timeval {
	return Timeval{
		Second: me.Second + 1,
		Micro:  0,
	}
}

func (me Timeval) Down() Timeval {
	return Timeval{
		Second: me.Second,
		Micro:  0,
	}
}

func (me Timeval) Time64() Time64 {
	return MakeTime64(me.Second, 1000*me.Micro)
}

func (me Timeval) Timespec() Timespec {
	return Timespec{
		Second: me.Second,
		Nano:   1000 * me.Micro,
	}
}

func (me *Timeval) String() string {
	return time.Unix(int64(me.Second), int64(1000*me.Micro)).Format(TimeFormat)
}

func (me *Timeval) Load(t Time64) {
	me.Second, me.Micro = t.Split()

	me.Micro /= 1000
}

func (me Timeval) Compare(v Timeval) (int, Timeval /*diff*/) {
	a := me.Time64()
	b := me.Time64()

	cmp, diff := a.Compare(b)

	return cmp, diff.Timeval()
}

func (me Timeval) inRange(a, b Timeval) bool {
	if cmp, _ := me.Compare(a); cmp < 0 { // me < a
		return false
	}

	if cmp, _ := me.Compare(b); cmp > 0 { // me > b
		return false
	}

	return true
}

func (me Timeval) InRange(a, b Timeval) bool {
	cmp, _ := a.Compare(b)
	if cmp < 0 {
		// a < b
		return me.inRange(a, b)
	} else {
		// a >= b
		return me.inRange(b, a)
	}
}

func (me Timeval) InZone(z Timezone) bool {
	return me.inRange(z.Begin.Timeval(), z.End.Timeval())
}

func (me Timeval) Add(v Timeval) Timeval {
	t := me.Time64() + v.Time64()

	return t.Timeval()
}

/******************************************************************************/
type Timezone32 struct {
	Begin Time32 `json:"begin"`
	End   Time32 `json:"end"`
}

func (me Timezone32) String() string {
	return fmt.Sprintf("begin(%s) end(%s)", me.Begin, me.End)
}

func (me Timezone32) IsGood() bool {
	return me.Begin.IsGood() || me.End.IsGood()
}

func (me *Timezone32) Zero() {
	me.Begin = 0
	me.End = 0
}

func (me Timezone32) Timezone() Timezone {
	return Timezone{
		Begin: Timespec{Second: me.Begin},
		End:   Timespec{Second: me.End},
	}
}

func (me Timezone32) Include(z Timezone32) bool {
	// |--------- me ---------|
	//      |----- z -----|
	return me.Begin <= z.Begin && me.End >= z.End
}

func (me Timezone32) Intersect(v Timezone32) Timezone32 {
	if v.Begin.InZone(me) {
		if v.End.InZone(me) {
			// |--------- me ---------|
			//     |----- v -----|
			return v
		} else {
			// |--------- me ---------|
			//               |----- v -----|
			return Timezone32{
				Begin: v.Begin,
				End:   me.End,
			}
		}
	} else {
		if v.End.InZone(me) {
			//     |--------- me ---------|
			// |----- v -----|
			return Timezone32{
				Begin: me.Begin,
				End:   v.End,
			}
		} else {
			//                  |--------- me ---------|
			// |----- v -----|              or              |----- v -----|
			return Timezone32{}
		}
	}
}

func (me Timezone32) Compare(v Timezone32) int {
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

/******************************************************************************/

type Timezone64 = Timezone

type Timezone struct {
	Begin Timespec `json:"begin"`
	End   Timespec `json:"end"`
}

func (me Timezone) String() string {
	return fmt.Sprintf("begin(%s %d) end(%s %d)",
		me.Begin, me.Begin.Nano,
		me.End, me.End.Nano)
}

func (me Timezone) IsGood() bool {
	return me.Begin.IsGood() && me.End.IsGood()
}

func (me Timezone) Timezone32() Timezone32 {
	return Timezone32{
		Begin: me.Begin.Second,
		End:   me.End.Second,
	}
}

func (me *Timezone) Zero() {
	me.Begin.Zero()
	me.End.Zero()
}

func (me Timezone) Include(z Timezone) bool {
	// |--------- me ---------|
	//      |----- z -----|
	a, _ := me.Begin.Compare(z.Begin)
	b, _ := me.End.Compare(z.End)

	return a <= 0 && b >= 0
}

func (me Timezone) Match(v Timezone) bool {
	return v.Begin.InZone(me) || v.End.InZone(me)
}

func (me Timezone) Intersect(v Timezone) Timezone {
	if v.Begin.InZone(me) {
		if v.End.InZone(me) {
			// |--------- me ---------|
			//     |----- v -----|
			return v
		} else {
			// |--------- me ---------|
			//               |----- v -----|
			return Timezone{
				Begin: v.Begin,
				End:   me.End,
			}
		}
	} else {
		if v.End.InZone(me) {
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

func (me Timezone) Compare(v Timezone) int {
	if cmp, _ := me.End.Compare(v.Begin); cmp < 0 {
		// |--------- me ---------|
		//                            |----- v -----|
		return -1
	} else if cmp, _ := me.Begin.Compare(v.End); cmp > 0 {
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

/******************************************************************************/

type TimeTask struct {
	Begin string `json:"begin"`
	End   string `json:"end"`
	Used  int
}
