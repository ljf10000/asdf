package asdf

type WorkerProgress struct {
	Begin   int    `json:"begin"`
	End     int    `json:"end"`
	Current int    `json:"current"`
	Percent string `json:"percent"`
}

func (me *WorkerProgress) Calc() {
	if (1 + me.Current) < me.End {
		count := me.End - me.Begin
		diff := 1 + me.End - me.Current

		percent := 100 * diff / count

		me.Percent = Itoa(percent) + "%"
	} else {
		me.Percent = "100%"
	}
}
