package asdf

import (
	"encoding/hex"
	"sync/atomic"
)

const (
	SizeofID64 = SizeofInt64

	LengthofID64String = 2 * SizeofID64
)

var ErrBadId64StringLengh error = ErrSprintf("bad id64 string length")

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
