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
func NewStimer(name string, count int, data interface{}) *Stimer {
	timer := &Stimer{
		name:  name,
		data:  data,
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
	data  interface{}

	slots []list.List
}

func (me *Stimer) index(idx int) int {
	return idx % len(me.slots)
}

func (me *Stimer) slot(idx int) *list.List {
	return &me.slots[me.index(idx)]
}

func (me *Stimer) Data() interface{} {
	return me.data
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
	if nil == node.GetElm() {
		idx := me.index(me.cur + after)
		slot := me.slot(idx)

		elm := slot.PushFront(node)

		node.SetElm(elm)
		node.SetSlot(idx)

		me.count++
	}
}

func (me *Stimer) Remove(node IStimerNode) {
	elm := node.GetElm()
	if nil != elm {
		slot := me.slot(node.GetSlot())
		slot.Remove(elm)

		node.SetElm(nil)

		me.count--
	}
}

func (me *Stimer) Trigger() {
	slot := me.slot(me.Current())

	elm := slot.Front()
	for nil != elm {
		node, _ := elm.Value.(IStimerNode)

		err := node.Handle(me)
		if nil != err {
			Log.Error("timer[%s] handle node[%s] error:%s", me.name, node, err)
		}

		me.Remove(node)
		elm = slot.Front()
	}

	me.cur++
	me.ticks++
}

func (me *Stimer) TriggerAll() {
	count := len(me.slots)

	for i := 0; i < count; i++ {
		me.Trigger()
	}
}

func (me *Stimer) Stop() {
	count := len(me.slots)

	for i := 0; i < count; i++ {
		slot := me.slot(i)

		elm := slot.Front()
		for nil != elm {
			node, _ := elm.Value.(IStimerNode)

			me.Remove(node)
			elm = slot.Front()
		}
	}
}
