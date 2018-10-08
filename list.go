package asdf

/******************************************************************************/
const SizeofListNode = 2 * SizeofPointer

type ListNode struct {
	prev *ListNode
	next *ListNode
}

func (me *ListNode) Init() {
	me.prev = nil
	me.next = nil
}

func (me *ListNode) add(prev *ListNode, next *ListNode) {
	next.prev = me
	me.next = next
	me.prev = prev
	prev.next = me
}

func (me *ListNode) del(prev *ListNode, next *ListNode) {
	next.prev = prev
	prev.next = next
}

/******************************************************************************/

type List struct {
	count int
	list  ListNode
}

func (me *List) Init() {
	me.list.prev = &me.list
	me.list.next = &me.list

	me.count = 0
}

func (me *List) Count() int {
	return me.count
}

func (me *List) Add(obj *ListNode) {
	obj.add(&me.list, me.list.next)

	me.count++
}

func (me *List) AddTail(obj *ListNode) {
	obj.add(me.list.prev, &me.list)

	me.count++
}

func (me *List) Del(obj *ListNode) {
	obj.del(obj.prev, obj.next)
	obj.Init()

	me.count--
}

func (me *List) First() *ListNode {
	return me.list.next
}

func (me *List) Tail() *ListNode {
	return me.list.prev
}

func (me *List) IsLast(obj *ListNode) bool {
	return obj.next == &me.list
}

func (me *List) IsFirst(obj *ListNode) bool {
	return obj.prev == &me.list
}

func (me *List) IsEmpty() bool {
	return 0 == me.count
}

func (me *List) Foreach(handle func(node *ListNode) error) error {
	for node := me.list.next; node != &me.list; node = node.next {
		err := handle(node)
		if nil != err {
			return err
		}
	}

	return nil
}

func (me *List) ForeachR(handle func(node *ListNode) error) error {
	for node := me.list.prev; node != &me.list; node = node.prev {
		err := handle(node)
		if nil != err {
			return err
		}
	}

	return nil
}

/******************************************************************************/

type HListNode struct {
	next  *HListNode
	pprev **HListNode
}

func (me *HListNode) Init() {
	me.pprev = nil
	me.next = nil
}

func (me *HListNode) IsHashed() bool {
	return nil != me.pprev
}

func (me *HListNode) del() {
	if !me.IsHashed() {
		next := me.next
		pprev := me.pprev
		*pprev = next

		if nil != next {
			next.pprev = pprev
		}

		me.Init()
	}
}

/******************************************************************************/

type hlistBucket struct {
	first *HListNode
}

func (me *hlistBucket) isEmpty() bool {
	return nil == me.first
}

func (me *hlistBucket) add(node *HListNode) {
	first := me.first
	node.next = first
	if nil != first {
		first.pprev = &node.next
	}
	me.first = node
	node.pprev = &me.first
}

func (me *hlistBucket) del(node *HListNode) {
	node.del()
}

func (me *hlistBucket) foreach(handle func(node *HListNode) error) error {
	for node := me.first; nil != node; node = node.next {
		err := handle(node)
		if nil != err {
			return err
		}
	}

	return nil
}

func (me *hlistBucket) find(eq HListEq) (*HListNode, bool) {
	for node := me.first; nil != node; node = node.next {
		if eq(node) {
			return node, true
		}
	}

	return nil, false
}

/******************************************************************************/

type HListIndex func() int
type HListEq func(node *HListNode) bool

type HList struct {
	buckets []hlistBucket
	count   int // all node count
}

func (me *HList) Count() int {
	return me.count
}

func (me *HList) Add(node *HListNode, index HListIndex, eq HListEq) (*HListNode, bool) {
	count := len(me.buckets)
	idx := index() % count
	bucket := &me.buckets[idx]

	obj, ok := bucket.find(eq)
	if ok {
		return obj, true
	}

	bucket.add(node)
	me.count++

	return node, false
}

func (me *HList) Del(node *HListNode, index HListIndex) {
	count := len(me.buckets)
	idx := index() % count
	bucket := &me.buckets[idx]

	if !bucket.isEmpty() {
		bucket.del(node)

		me.count--
	}
}
