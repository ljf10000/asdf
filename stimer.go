package asdf

import (
	"container/list"
)

type IStimerNode interface {
	String() string

	GetSlot() int
	SetSlot(slot int)

	GetElm() *list.Element
	SetElm(elm *list.Element)

	Handle(timer *Stimer) error
}

/******************************************************************************/
func NewStimer(name string, count int) *Stimer {
	timer := &Stimer{
		name:  name,
		slots: make([]list.List, count),
	}

	for i := 0; i < count; i++ {
		timer.slots[i].Init()
	}

	return timer
}

// simple timer
type Stimer struct {
	name  string
	ticks uint64
	cur   int
	count int

	slots []list.List
}

func (me *Stimer) index(idx int) int {
	return idx % len(me.slots)
}

func (me *Stimer) slot(idx int) *list.List {
	return &me.slots[me.index(idx)]
}

func (me *Stimer) Window() int {
	return len(me.slots)
}

func (me *Stimer) Ticks() uint64 {
	return me.ticks
}

func (me *Stimer) Len() int {
	return me.count
}

func (me *Stimer) SlotLen(idx int) int {
	return me.slot(idx).Len()
}

func (me *Stimer) Current() int {
	return me.index(me.cur)
}

func (me *Stimer) Insert(node IStimerNode, after int) {
	idx := me.index(me.cur + after)
	slot := me.slot(idx)

	elm := slot.PushFront(node)

	node.SetElm(elm)
	node.SetSlot(me.Current())

	me.count++
}

func (me *Stimer) Remove(node IStimerNode) {
	slot := me.slot(node.GetSlot())

	slot.Remove(node.GetElm())

	me.count--
}

func (me *Stimer) Trigger() {
	slot := me.slot(me.Current())

	elm := slot.Front()
	for nil != elm {
		v := slot.Remove(elm)

		node, _ := v.(IStimerNode)

		err := node.Handle(me)
		if nil != err {
			Log.Error("timer[%s] handle node[%s] error:%s", me.name, node, err)
		}

		elm = slot.Front()
	}

	me.cur++
	me.ticks++
}
