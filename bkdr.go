package asdf

const BkdrFactor = 37

type Bkdr uint32

type IBkdr interface {
	Bkdr(buf []byte) Bkdr
}

type DeftBkdr struct{}

func GenBkdr(buf []byte) Bkdr {
	bkdr := uint64(BkdrFactor)

	for _, b := range buf {
		bkdr = bkdr*BkdrFactor + uint64(b)
	}

	return Bkdr(bkdr)
}

func (me *DeftBkdr) Bkdr(buf []byte) Bkdr {
	return GenBkdr(buf)
}

var DeftBkdrer = &DeftBkdr{}
