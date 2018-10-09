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

func (me *ListNode) InList() bool {
	return nil != me.next || nil != me.prev
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

	me.Init()
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

func (me *List) Add(node *ListNode) {
	if !node.InList() {
		node.add(&me.list, me.list.next)

		me.count++
	}
}

func (me *List) AddTail(node *ListNode) {
	if !node.InList() {
		node.add(me.list.prev, &me.list)

		me.count++
	}
}

func (me *List) Del(node *ListNode) {
	if me.count > 0 && node.InList() {
		node.del(node.prev, node.next)

		me.count--
	}
}

func (me *List) First() *ListNode {
	if 0 == me.count {
		return nil
	} else {
		return me.list.next
	}
}

func (me *List) Tail() *ListNode {
	if 0 == me.count {
		return nil
	} else {
		return me.list.prev
	}
}

func (me *List) IsLast(node *ListNode) bool {
	return me.count > 0 && node.next == &me.list
}

func (me *List) IsFirst(node *ListNode) bool {
	return me.count > 0 && node.prev == &me.list
}

func (me *List) IsEmpty() bool {
	return 0 == me.count
}

func (me *List) Foreach(handle func(node *ListNode) error) error {
	if me.count > 0 {
		for node := me.list.next; node != &me.list; node = node.next {
			err := handle(node)
			if nil != err {
				return err
			}
		}
	}

	return nil
}

func (me *List) ForeachR(handle func(node *ListNode) error) error {
	if me.count > 0 {
		for node := me.list.prev; node != &me.list; node = node.prev {
			err := handle(node)
			if nil != err {
				return err
			}
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

func (me *HListNode) InHash() bool {
	return nil != me.pprev
}

func (me *HListNode) del() {
	if me.InHash() {
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

func (me *hlistBucket) clean() {

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

func (me *hlistBucket) find(eq HListEq) *HListNode {
	if me.isEmpty() {
		return nil
	}

	for node := me.first; nil != node; node = node.next {
		if eq(node) {
			return node
		}
	}

	return nil
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

func (me *HList) Init(count int) {
	me.buckets = make([]hlistBucket, count)
}

func (me *HList) bucket(index HListIndex) *hlistBucket {
	count := len(me.buckets)
	idx := index() % count

	return &me.buckets[idx]
}

func (me *HList) Add(node *HListNode, index HListIndex, eq HListEq) (*HListNode, bool) {
	bucket := me.bucket(index)

	obj := bucket.find(eq)
	if nil != obj {
		return obj, true
	}

	bucket.add(node)
	me.count++

	return node, false
}

func (me *HList) Del(node *HListNode, index HListIndex) {
	bucket := me.bucket(index)

	if !bucket.isEmpty() {
		bucket.del(node)

		me.count--
	}
}

func (me *HList) Get(index HListIndex, eq HListEq) *HListNode {
	bucket := me.bucket(index)

	return bucket.find(eq)
}
