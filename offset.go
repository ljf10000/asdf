package asdf

import (
	"encoding/binary"
)

type Offsetof struct {
	Offset uint32
	Len    uint32
}

func (me *Offsetof) Set(Offset, Len int) {
	me.Offset = uint32(Offset)
	me.Len = uint32(Len)
}

func (me *Offsetof) Size() int {
	return 2 * SizeofInt32
}

func (me *Offsetof) ToBinary(bin []byte) error {
	if len(bin) < me.Size() {
		return ErrTooShortBuffer
	}

	binary.BigEndian.PutUint32(bin[0:], me.Offset)
	binary.BigEndian.PutUint32(bin[4:], me.Len)

	return nil
}

func (me *Offsetof) FromBinary(bin []byte) error {
	if len(bin) < me.Size() {
		return ErrTooShortBuffer
	}

	me.Offset = binary.BigEndian.Uint32(bin[0:])
	me.Len = binary.BigEndian.Uint32(bin[4:])

	return nil
}
