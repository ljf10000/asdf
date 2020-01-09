package asdf

const (
	FrequencyLow    Frequency = 0
	FrequencyNormal Frequency = 1
	FrequencyHigh   Frequency = 2
	FrequencyEnd    Frequency = 4
)

type Frequency byte

var frequencies = EnumMapper{
	Type: "asdf.Frequency",
	Names: []string{
		FrequencyLow:    "low",
		FrequencyNormal: "normal",
		FrequencyHigh:   "high",
	},
}

func (me Frequency) IsGood() bool {
	return me < FrequencyEnd
}

func (me Frequency) String() string {
	if me.IsGood() {
		return frequencies.Name(int(me))
	} else {
		return Unknow
	}
}

func (me *Frequency) FromString(s string) error {
	idx, err := frequencies.Index(s)
	if nil == err {
		*me = Frequency(idx)
	}

	return err
}
