package asdf

import (
	"unsafe"
)

const (
	INVALID_STIMER_SLOT uint32 = 0xffffffff

	SizeofUnsafeStimerNode = SizeofListNode + SizeofPointer + 2*SizeofInt32
)

var scUnsafeStimerNode = NewSizeChecker("UnsafeStimerNode", unsafe.Sizeof(UnsafeStimerNode{}), SizeofUnsafeStimerNode)

type UnsafeStimerNode struct {
	Handle func(node *UnsafeStimerNode) error

	node   ListNode
	slot   uint32
	expire uint32
}

func (me *UnsafeStimerNode) Init() {
	me.node.Init()
}

func (me *UnsafeStimerNode) Expire() int {
	return int(me.expire)
}

/******************************************************************************/
func NewUnsafeStimer(name string, count int) *UnsafeStimer {
	timer := &UnsafeStimer{
		name:  name,
		slots: make([]List, count),
	}

	for i := 0; i < count; i++ {
		slot := timer.slot(i)
		slot.Init()
	}

	return timer
}

// simple timer
type UnsafeStimer struct {
	name  string
	ticks uint64
	cur   int
	count int

	slots []List
}

func (me *UnsafeStimer) index(idx int) int {
	return idx % me.Window()
}

func (me *UnsafeStimer) slot(idx int) *List {
	return &me.slots[me.index(idx)]
}

func (me *UnsafeStimer) Window() int {
	return len(me.slots)
}

func (me *UnsafeStimer) Ticks() uint64 {
	return me.ticks
}

func (me *UnsafeStimer) Len() int {
	return me.count
}

func (me *UnsafeStimer) Current() int {
	return me.index(me.cur)
}

func (me *UnsafeStimer) Insert(node *UnsafeStimerNode, after int) {
	expire := me.cur + after
	idx := me.index(expire)
	slot := me.slot(idx)

	slot.InsertHead(&node.node)
	node.slot = uint32(idx)
	node.expire = uint32(expire)
	me.count++
}

func (me *UnsafeStimer) Remove(node *UnsafeStimerNode) {
	idx := int(node.slot)
	slot := me.slot(idx)
	slot.Remove(&node.node)
	node.slot = INVALID_STIMER_SLOT

	me.count--
}

func _ListNode_to_UnsafeStimerNode(node *ListNode) *UnsafeStimerNode {
	return (*UnsafeStimerNode)(unsafe.Pointer(node))
}

func (me *UnsafeStimer) Trigger() {
	slot := me.slot(me.Current())

	pos := slot.First()
	for nil != pos {
		node := _ListNode_to_UnsafeStimerNode(pos)

		err := node.Handle(node)
		if nil != err {
			Log.Error("timer[%s] handle node[%s] error:%s", me.name, node, err)
		}

		me.Remove(node)
		pos = slot.First()
	}

	me.cur++
	me.ticks++
}

func (me *UnsafeStimer) TriggerAll() {
	count := len(me.slots)

	for i := 0; i < count; i++ {
		me.Trigger()
	}
}

func (me *UnsafeStimer) Stop() {
	count := me.Window()

	for i := 0; i < count; i++ {
		slot := me.slot(i)

		pos := slot.First()
		for nil != pos {
			node := _ListNode_to_UnsafeStimerNode(pos)

			me.Remove(node)
			pos = slot.First()
		}
	}
}
