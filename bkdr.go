package asdf

import (
	"io"
)

const BkdrFactor = 37

type Bkdr uint32

type IBkdrWriter interface {
	Write(buf []byte) (int, error)
}

type IBkdr interface {
	io.Writer

	Bkdr(buf []byte) Bkdr
}

type DeftBkdr struct {
	bkdr Bkdr
}

func (me *DeftBkdr) Bkdr(buf []byte) Bkdr {
	me.Write(buf)

	return me.bkdr
}

func (me *DeftBkdr) Write(buf []byte) (int, error) {
	bkdr := uint64(BkdrFactor)

	for _, b := range buf {
		bkdr = bkdr*BkdrFactor + uint64(b)
	}

	me.bkdr = Bkdr(bkdr)

	return len(buf), nil
}

var DeftBkdrer = &DeftBkdr{}
