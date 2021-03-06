package asdf

import (
	"encoding/binary"
	"math/rand"
	"time"
)

var RandSeed = NewRandSeed()

func NewRandSeed() *rand.Rand {
	return rand.New(rand.NewSource(time.Now().UnixNano()))
}

type RandBuffer []byte

func GenRandBuffer(buf []byte) {
	RandSeed.Read(buf)
}

func (me RandBuffer) Rand() {
	GenRandBuffer([]byte(me))
}

func (me RandBuffer) RandTime() {
	count := len(me)
	offset := 0

	if count > 8 {
		offset = 8
		binary.BigEndian.PutUint64(me[0:8], uint64(time.Now().UnixNano()))
	} else if count > 4 {
		offset = 4
		binary.BigEndian.PutUint32(me[0:4], uint32(time.Now().Unix()))
	} else {
		return
	}

	for i := 0; i < count-offset; i++ {
		me[i+offset] = byte(RandSeed.Uint32())
	}
}
