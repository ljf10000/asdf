package asdf

const (
	EvTick  Event = 0
	EvExit  Event = 1
	EvQuit  Event = 2
	EvStop  Event = 3
	EvIoErr Event = 4

	EvEnd Event = 5
)

type Event byte

var evTypes = &EnumMapper{
	Type: "asdf.Event",
	Names: []string{
		EvStop:  "stop",
		EvTick:  "tick",
		EvExit:  "exit",
		EvQuit:  "quit",
		EvIoErr: "ioerr",
	},
}

func (me Event) IsGood() bool {
	return evTypes.IsGoodIndex(int(me))
}

func (me Event) String() string {
	return evTypes.Name(int(me))
}

func (me *Event) FromString(s string) error {
	idx, err := evTypes.Index(s)
	if nil == err {
		*me = Event(idx)
	}

	return err
}
