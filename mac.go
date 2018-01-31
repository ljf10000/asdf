package asdf

import (
	"encoding/hex"
	"fmt"
)

const (
	MacSize    = 6
	MacStringS = 12 // AABBCCDDEEFF
	MacStringM = 14 // AABB-CCDD-EEFF
	MacStringL = 17 // AA:BB:CC:DD:EE:FF or AA-BB-CC-DD-EE-FF

	MacSepWindows = '-'
	MacSepUnix    = ':'
)

type MAC [MacSize]byte
type Mac []byte

func (me MAC) Mac() Mac {
	return Mac(me[:])
}

func (me MAC) String() string {
	return me.Mac().String()
}

func (me Mac) IsGood() bool {
	return MacSize == len(me) && !Slice(me).IsZero() && !Slice(me).IsFull()
}

func (me Mac) Slice() []byte {
	return me
}

func (me Mac) Eq(it interface{}) bool {
	return Slice(me).Eq(it)
}

func (me Mac) ToStringL(ifs byte) string {
	return fmt.Sprintf("%.2x%c%.2x%c%.2x%c%.2x%c%.2x%c%.2x",
		me[0], ifs,
		me[1], ifs,
		me[2], ifs,
		me[3], ifs,
		me[4], ifs,
		me[5])
}

func (me Mac) ToStringLU() string {
	return me.ToStringL(MacSepUnix)
}

func (me Mac) ToStringLW() string {
	return me.ToStringL(MacSepWindows)
}

func (me Mac) ToStringM(ifs byte) string {
	return fmt.Sprintf("%.2x%.2x%c%.2x%.2x%c%.2x%.2x",
		me[0], me[1], ifs,
		me[2], me[3], ifs,
		me[4], me[5])
}

func (me Mac) ToStringMU() string {
	return me.ToStringM(MacSepUnix)
}

func (me Mac) ToStringMW() string {
	return me.ToStringM(MacSepWindows)
}

func (me Mac) ToStringS() string {
	return fmt.Sprintf("%.2x%.2x%.2x%.2x%.2x%.2x",
		me[0], me[1],
		me[2], me[3],
		me[4], me[5])
}

func (me Mac) String() string {
	return me.ToStringLU()
}

func (me Mac) FromString(s string) error {
	b := []byte(s)

	switch len(s) {
	case MacStringL: // AA:BB:CC:DD:EE:FF or AA-BB-CC-DD-EE-FF
		ifs := b[2]

		if (ifs != MacSepUnix && ifs != MacSepWindows) ||
			ifs != b[5] ||
			ifs != b[8] ||
			ifs != b[11] ||
			ifs != b[14] {

			return Error
		}

		for i := 0; i < 6; i++ {
			if _, err := hex.Decode(me[i:], b[3*i:3*i+2]); nil != err {
				return err
			}
		}

		return nil
	case MacStringM: // AABB-CCDD-EEFF or AABB:CCDD:EEFF
		ifs := b[4]

		if (ifs != MacSepUnix && ifs != MacSepWindows) ||
			ifs != b[9] {
			return Error
		}

		for i := 0; i < 3; i++ {
			if _, err := hex.Decode(me[2*i:], b[5*i:5*i+2]); nil != err {
				return err
			}
			if _, err := hex.Decode(me[2*i+1:], b[5*i+2:5*i+4]); nil != err {
				return err
			}
		}

		return nil
	case MacStringS: // AABBCCDDEEFF
		_, err := hex.Decode(me[:], b)

		return err
	default:
		return ErrBadMac
	}
}

type MacString string

func (me MacString) IsGood() bool {
	mac := MAC{}

	if err := Mac(mac[:]).FromString(string(me)); nil != err {
		return false
	}

	return true
}

func (me MacString) ToBinary() MAC {
	mac := MAC{}

	Mac(mac[:]).FromString(string(me))

	return mac
}

func (me MacString) ToString() string {
	return string(me)
}
