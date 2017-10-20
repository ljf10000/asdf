package asdf

import (
	"encoding/binary"
	"encoding/hex"
	"sync/atomic"
	"time"
)

const suidSize = 16

type Suid [suidSize]byte

func (me *Suid) Seq() uint32 {
	return binary.BigEndian.Uint32((*me)[12:])
}

func (me *Suid) ToString() string {
	return hex.EncodeToString((*me)[:])
}

func (me *Suid) FromString(s string) error {
	if len(s) != 2*suidSize {
		return ErrBadLen
	} else if buf, err := hex.DecodeString(s); nil != err {
		return err
	} else {
		copy((*me)[:], buf)

		return nil
	}
}

type ISuID interface {
	Generate() Suid
}

type suidGenerator struct {
	name    string
	unix    uint32
	rand    [8]byte
	address *uint32
}

func (me suidGenerator) Generate() Suid {
	suid := Suid{}

	binary.BigEndian.PutUint32(suid[0:4], me.unix)
	copy(suid[4:12], me.rand[:])
	binary.BigEndian.PutUint32(suid[12:], atomic.AddUint32(me.address, 1))

	return suid
}

func SuidGenerator(name string) ISuID {
	seqAddress := (*uint32)(nil)

	suidLock.Handle(func() {
		address, ok := suidMap[name]
		if !ok {
			seq := uint32(RandSeed.Int31() & 0xffff)
			address = &seq
			suidMap[name] = address
		}

		seqAddress = address
	})

	generator := &suidGenerator{
		name:    name,
		address: seqAddress,
		unix:    uint32(time.Now().Unix()),
	}

	RandBuffer(generator.rand[:]).Rand()

	Log.Debug("new suid generator:%+v", generator)

	return generator
}

var suidLock = &AccessLock{}
var suidMap = make(map[string]*uint32)
