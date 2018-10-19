package asdf

import (
	"sync/atomic"
	"unsafe"
)

const (
	SizeofStatCounter   = SizeofCacheLine
	CountofStatCouonter = SizeofStatCounter/SizeofInt64 - 1
)

// align of cache line
type StatCounter struct {
	v uint64
	_ [CountofStatCouonter]uint64
}

func (me *StatCounter) Add(v int) uint64 {
	return atomic.AddUint64(&me.v, uint64(v))
}

func (me *StatCounter) Get() uint64 {
	return atomic.LoadUint64(&me.v)
}

func (me *StatCounter) Set(v uint64) {
	atomic.StoreUint64(&me.v, v)
}

/******************************************************************************/
const (
	SizeofStatePage = 4 * SizeofK
	CountofStatPage = SizeofStatePage/SizeofStatCounter - 1
)

type DeftStat struct {
	start StatCounter

	counter [CountofStatPage]StatCounter
}

func (me *DeftStat) Start() {
	me.start.Set(0)

	for i := 0; i < CountofStatPage; i++ {
		me.counter[i].Set(0)
	}

	me.start.Set(1)
}

func (me *DeftStat) Stop() {
	me.start.Set(0)
}

func (me *DeftStat) IsStart() bool {
	return 1 == me.start.Get()
}

func (me *DeftStat) IsGoodIndex(idx int) bool {
	return idx >= 0 && idx < CountofStatPage
}

func (me *DeftStat) Add(idx int, count int) {
	if me.IsGoodIndex(idx) && me.IsStart() {
		me.counter[idx].Add(count)
	}
}

func (me *DeftStat) Get(idx int) uint64 {
	if me.IsGoodIndex(idx) && me.IsStart() {
		return me.counter[idx].Get()
	} else {
		return 0
	}
}

func (me *DeftStat) Size() int {
	return SizeofStatePage
}

func (me *DeftStat) Slice() []byte {
	return MakeSlice(unsafe.Pointer(me), SizeofStatePage, SizeofStatePage)
}
