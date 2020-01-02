package asdf

const (
	EvTick       Event = 0
	EvExit       Event = 1
	EvQuit       Event = 2
	EvStop       Event = 3
	EvClose      Event = 4
	EvReadEof    Event = 5
	EvReadError  Event = 6
	EvWriteEof   Event = 7
	EvWriteError Event = 8

	EvEnd Event = 9
)

type Event byte

var evTypes = &EnumMapper{
	Type: "asdf.Event",
	Names: []string{
		EvTick:       "tick",
		EvExit:       "exit",
		EvQuit:       "quit",
		EvStop:       "stop",
		EvClose:      "close",
		EvReadEof:    "reof",
		EvReadError:  "rerr",
		EvWriteEof:   "weof",
		EvWriteError: "werr",
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
