package asdf

type WorkerProgress struct {
	Begin   int    `json:"begin"`
	End     int    `json:"end"`
	Count   int    `json:"count"`
	Current int    `json:"current"`
	Percent string `json:"percent"`
}

func (me *WorkerProgress) String() string {
	return "begin:" + Itoa(me.Begin) +
		", end:" + Itoa(me.End) +
		", current:" + Itoa(me.Current) +
		", percent:" + me.Percent
}

func (me *WorkerProgress) Calc() {
	me.Count = me.End - me.Begin

	if (1 + me.Current) < me.End {
		diff := 1 + me.Current - me.Begin

		percent := 100 * diff / me.Count

		me.Percent = Itoa(percent) + "%"
	} else {
		me.Percent = "100%"
	}
}
