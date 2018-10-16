package asdf

import (
	"sync/atomic"
)

type StatCounter struct {
	v uint64
	_ [7]uint64
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
