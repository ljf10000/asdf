package asdf

import (
	"encoding/hex"
	"sync/atomic"
)

const (
	SizeofID32 = SizeofInt32
	SizeofID64 = SizeofInt64

	LengthofID32String = 2 * SizeofID32
	LengthofID64String = 2 * SizeofID64
)

var ErrBadId32StringLengh error = ErrSprintf("bad id32 string length")
var ErrBadId64StringLengh error = ErrSprintf("bad id64 string length")

type ID32 uint32

func (me *ID32) Add(v uint32) ID32 {
	return ID32(atomic.AddUint32((*uint32)(me), v))
}

func (me ID32) String() string {
	buf := [SizeofID32]byte{}
	dst := [LengthofID32String]byte{}

	Hton32(buf[:], uint32(me))

	hex.Encode(dst[:], buf[:])

	return string(dst[:])
}

func (me *ID32) FromString(s string) error {
	if LengthofID32String != len(s) {
		return ErrBadId32StringLengh
	}

	buf := [SizeofID32]byte{}

	if _, err := hex.Decode(buf[:], []byte(s)); nil != err {
		return err
	}

	*me = ID32(Ntoh32(buf[:]))

	return nil
}

type ID64 uint64

func (me *ID64) Add(v uint64) ID64 {
	return ID64(atomic.AddUint64((*uint64)(me), v))
}

func (me ID64) String() string {
	buf := [SizeofID64]byte{}
	dst := [LengthofID64String]byte{}

	Hton64(buf[:], uint64(me))

	hex.Encode(dst[:], buf[:])

	return string(dst[:])
}

func (me *ID64) FromString(s string) error {
	if LengthofID64String != len(s) {
		return ErrBadId64StringLengh
	}

	buf := [SizeofID64]byte{}

	if _, err := hex.Decode(buf[:], []byte(s)); nil != err {
		return err
	}

	*me = ID64(Ntoh64(buf[:]))

	return nil
}
