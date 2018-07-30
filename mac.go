package asdf

import (
	"encoding/hex"
	"fmt"
	"unsafe"
)

const (
	MacSize    = 6
	MacStringS = 12 // AABBCCDDEEFF
	MacStringM = 14 // AABB-CCDD-EEFF
	MacStringL = 17 // AA:BB:CC:DD:EE:FF or AA-BB-CC-DD-EE-FF

	MacSepWindows = '-'
	MacSepUnix    = ':'

	MAC_F_MULTICASE   = 0x01
	MAC_F_LOCAL_ADMIN = 0x02

	_uptr_4 = uintptr(4)
)

var ZERO_MAC = Mac{0, 0, 0, 0, 0, 0}
var FULL_MAC = Mac{0xff, 0xff, 0xff, 0xff, 0xff, 0xff}

type Mac []byte

func (me Mac) IsUnicast() bool {
	return 0 == (me[0] & MAC_F_MULTICASE)
}

func (me Mac) IsMulticast() bool {
	return MAC_F_MULTICASE == (me[0] & MAC_F_MULTICASE)
}

func (me Mac) IsUniversal() bool {
	return 0 == (me[0] & MAC_F_LOCAL_ADMIN)
}

func (me Mac) IsLocalAdmin() bool {
	return MAC_F_LOCAL_ADMIN == (me[0] & MAC_F_LOCAL_ADMIN)
}

func (me Mac) isSame(a uint32, b uint16) bool {
	p := SliceAddress(me)

	return *(*uint32)(unsafe.Pointer(p)) == a &&
		*(*uint16)(unsafe.Pointer(p + _uptr_4)) == b
}

func (me Mac) IsBroadcast() bool {
	return me.isSame(0xffffffff, 0xffff)
}

func (me Mac) IsZero() bool {
	return me.isSame(0, 0)
}

func (me Mac) IsGood() bool {
	return MacSize == len(me) && !me.IsZero()
}

func (me Mac) IsGoodUnicast() bool {
	return me.IsGood() && me.IsUnicast()
}

func (me Mac) Eq(mac Mac) bool {
	p := SliceAddress(mac)

	return MacSize == len(me) &&
		me.isSame(*(*uint32)(unsafe.Pointer(p)),
			*(*uint16)(unsafe.Pointer(p + _uptr_4)))
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
	mac := [MacSize]byte{}

	if err := Mac(mac[:]).FromString(string(me)); nil != err {
		return false
	}

	return true
}

func (me MacString) ToBinary() Mac {
	mac := [MacSize]byte{}

	Mac(mac[:]).FromString(string(me))

	return mac[:]
}

func (me MacString) ToString() string {
	return string(me)
}
