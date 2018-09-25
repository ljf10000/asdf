package asdf

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
