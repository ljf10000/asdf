package asdf

type IntZone struct {
	Min  int
	Max  int
	Deft int
}

func (me *IntZone) Repair(pv *int) {
	v := *pv

	if me.Min < v || v > me.Max {
		*pv = me.Deft
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
