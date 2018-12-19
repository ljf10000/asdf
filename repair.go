package asdf

const (
	INTZONE_F_IGNORE = 1
)

type IntZone struct {
	Name   string
	Flag   int
	Ignore int
	Min    int
	Max    int
	Deft   int
}

func (me *IntZone) CanIgnore() bool {
	return INTZONE_F_IGNORE == (INTZONE_F_IGNORE & me.Flag)
}

func (me *IntZone) IsIgnore(v int) bool {
	return me.CanIgnore() && v == me.Ignore
}

func (me *IntZone) CanRepair(v int) bool {
	return !me.IsIgnore(v) && (v < me.Min || v > me.Max)
}

func (me *IntZone) Repair(pv *int) {
	v := *pv

	if me.CanRepair(v) {
		*pv = me.Deft

		name := me.Name
		if Empty == name {
			name = Unknow
		}

		Log.Debug("intzone:%s min:%d max:%d deft:%d oldv:%d ==> newv:%d",
			name,
			me.Min, me.Max, me.Deft,
			v, *pv)
	}
}

func RepairInt(pv *int, begin, end, deft int) {
	v := *pv

	if v <= begin || v > end {
		*pv = deft
	}
}

func MinInt(a, b int) int {
	if a < b {
		return a
	} else {
		return b
	}
}

func MaxInt(a, b int) int {
	if a > b {
		return a
	} else {
		return b
	}
}
