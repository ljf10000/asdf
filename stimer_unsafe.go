package asdf

import (
	"unsafe"
)

const (
	INVALID_STIMER_SLOT uint32 = 0xffffffff

	SizeofUnsafeStimerNode = SizeofListNode + SizeofPointer + 2*SizeofInt32
)

var (
	zUnsafeStimerNode  = UnsafeStimerNode{}
	scUnsafeStimerNode = NewSizeChecker("UnsafeStimerNode", unsafe.Sizeof(UnsafeStimerNode{}), SizeofUnsafeStimerNode)
)

func unsafe_stimer_node_from(node unsafe.Pointer, offset uintptr) *UnsafeStimerNode {
	if nil != node {
		return (*UnsafeStimerNode)(GetObjByField(node, offset))
	} else {
		return nil
	}
}

func unsafe_stimer_node_from_listnode(node *ListNode) *UnsafeStimerNode {
	return unsafe_stimer_node_from(unsafe.Pointer(node), unsafe.Offsetof(zUnsafeStimerNode.node))
}

type UnsafeStimerHandle func(node *UnsafeStimerNode) error

type UnsafeStimerNode struct {
	node   ListNode
	slot   uint32
	expire uint32

	handle UnsafeStimerHandle
}

func (me *UnsafeStimerNode) Init(handle UnsafeStimerHandle) {
	me.node.Init()
	me.handle = handle
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

func (me *UnsafeStimer) Trigger() {
	slot := me.slot(me.Current())

	pos := slot.First()
	for nil != pos {
		node := unsafe_stimer_node_from_listnode(pos)

		Log.Debug("stimer[%s] trigger node[%p] pos[%p]", me.name, node, pos)

		err := node.handle(node)
		if nil != err {
			Log.Error("timer[%s] handle node[%p] error:%s", me.name, node, err)
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
			node := unsafe_stimer_node_from_listnode(pos)

			me.Remove(node)
			pos = slot.First()
		}
	}
}
