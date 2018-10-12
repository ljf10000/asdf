package asdf

type IntZone struct {
	Name string
	Min  int
	Max  int
	Deft int
}

func (me *IntZone) Repair(pv *int) {
	v := *pv

	if v < me.Min || v > me.Max {
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
