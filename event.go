package asdf

const (
	EvTick       Event = 0
	EvExit       Event = 1
	EvQuit       Event = 2
	EvStop       Event = 3
	EvReadEof    Event = 4
	EvReadError  Event = 5
	EvWriteEof   Event = 6
	EvWriteError Event = 7

	EvEnd Event = 8
)

type Event byte

var evTypes = &EnumMapper{
	Type: "asdf.Event",
	Names: []string{
		EvStop:       "stop",
		EvTick:       "tick",
		EvExit:       "exit",
		EvQuit:       "quit",
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
