package asdf

type WorkerProgress struct {
	Begin   int    `json:"begin"`
	End     int    `json:"end"`
	Count   int    `json:"count"`
	Handle  int    `json:"handle"`
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
	me.Handle = 1 + me.Current - me.Begin

	if (1 + me.Current) < me.End {
		percent := 100 * me.Handle / me.Count

		me.Percent = Itoa(percent) + "%"
	} else {
		me.Percent = "100%"
	}
}
