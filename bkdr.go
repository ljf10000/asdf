package asdf

const BkdrFactor Bkdr = 131

type Bkdr uint32

type IBkdr interface {
	Bkdr(buf []byte) Bkdr
}

type DeftBkdr struct{}

func GenBkdr(buf []byte) Bkdr {
	return MakeBkdr(0, buf)
}

func MakeBkdr(bkdr Bkdr, buf []byte) Bkdr {
	//Console.Info("make bkdr s=%s n=%d", string(buf), bkdr)
	for _, b := range buf {
		bkdr = bkdr*BkdrFactor + Bkdr(b)
		//Console.Info("make bkdr s=%s n=%d c=%d", string(buf), bkdr, b)
	}

	return bkdr
}

func (me *DeftBkdr) Bkdr(buf []byte) Bkdr {
	return GenBkdr(buf)
}

var DeftBkdrer = &DeftBkdr{}
