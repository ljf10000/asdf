package asdf

type WorkerProgress struct {
	Begin   int    `json:"begin"`
	End     int    `json:"end"`
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
	if (1 + me.Current) < me.End {
		count := me.End - me.Begin
		diff := 1 + me.Current - me.Begin

		percent := 100 * diff / count

		me.Percent = Itoa(percent) + "%"
	} else {
		me.Percent = "100%"
	}
}
