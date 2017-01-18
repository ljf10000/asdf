package asdf

type ChanMsg struct {
	Sender IObj
	Data   IObjType
}

func NewChanMsg(sender IObj, data IObjType) *ChanMsg {
	return &ChanMsg{
		Sender: sender,
		Data:   data,
	}
}
