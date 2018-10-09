package asdf

/******************************************************************************/

const SizeofListNode = 3 * SizeofPointer

type ListNode struct {
	list *List

	prev *ListNode
	next *ListNode
}

func (me *ListNode) Init() {
	me.list = nil
	me.prev = nil
	me.next = nil
}

func (me *ListNode) InList() bool {
	return nil != me.next || nil != me.prev
}

func (me *ListNode) inTheList(list *List) bool {
	return me.list == list
}

func (me *ListNode) InTheList(list *List) bool {
	return me.InList() && me.inTheList(list)
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

const SizeofList = SizeofListNode + SizeofInt64

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

func (me *List) InsertAfter(node *ListNode, obj *ListNode) error {
	if node.InList() {
		Log.Error("insert node:%p into list:%p error: have in list:%p", node, me, node.list)

		return ErrExist
	} else {
		node.add(obj, obj.next)
		node.list = me

		me.count++

		return nil
	}
}

func (me *List) InsertBefore(node *ListNode, obj *ListNode) error {
	if node.InList() {
		Log.Error("insert node:%p into list:%p error: have in list:%p", node, me, node.list)

		return ErrExist
	} else {
		node.add(obj.prev, obj)
		node.list = me

		me.count++

		return nil
	}
}

func (me *List) InsertHead(node *ListNode) error {
	if node.InList() {
		Log.Error("insert node:%p into list:%p error: have in list:%p", node, me, node.list)

		return ErrExist
	} else {
		node.add(&me.list, me.list.next)
		node.list = me

		me.count++

		return nil
	}
}

func (me *List) InsertTail(node *ListNode) error {
	if node.InList() {
		Log.Error("insert node:%p into list:%p error: have in list:%p", node, me, node.list)

		return ErrExist
	} else {
		node.add(me.list.prev, &me.list)
		node.list = me

		me.count++

		return nil
	}
}

func (me *List) Remove(node *ListNode) error {
	if me.count > 0 {
		if node.InList() {
			if node.inTheList(me) {
				node.del(node.prev, node.next)
				node.Init()

				me.count--

				return nil
			} else {
				Log.Error("remove node:%p from list:%p error: but in the list:%p",
					node,
					me,
					node.list)
			}
		} else {
			Log.Error("remove node:%p from list:%p error: not in list", node, me)
		}
	} else {
		Log.Error("remove node:%p from list:%p error: empty list", node, me)
	}

	return ErrNoExist
}

func (me *List) First() *ListNode {
	if 0 == me.count {
		return nil
	} else {
		return me.list.next
	}
}

func (me *List) Last() *ListNode {
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

type ListForeach func(node *ListNode) (error, bool)
type ListForeach2 func(a *ListNode, b *ListNode) (error, bool)

func (me *List) Foreach(handle ListForeach) (error, bool) {
	if me.count > 0 {
		for node := me.list.next; node != &me.list; node = node.next {
			if err, br := handle(node); br {
				return err, true
			}
		}
	}

	return nil, false
}

func (me *List) Foreach2(handle ListForeach2) (error, bool) {
	if me.count >= 2 {
		a := me.list.next
		b := a.next

		for b != &me.list {
			if err, br := handle(a, b); br {
				return err, true
			}

			a = b
			b = a.next
		}
	}

	return nil, false
}

func (me *List) SafeForeach(handle ListForeach) (error, bool) {
	if me.count > 0 {
		node := me.list.next
		next := node.next

		for node != &me.list {
			if err, br := handle(node); br {
				return err, true
			}

			node = next
			next = node.next
		}
	}

	return nil, false
}

func (me *List) ForeachR(handle ListForeach) (error, bool) {
	if me.count > 0 {
		for node := me.list.prev; node != &me.list; node = node.prev {
			if err, br := handle(node); br {
				return err, true
			}
		}
	}

	return nil, false
}

func (me *List) SafeForeachR(handle ListForeach) (error, bool) {
	if me.count > 0 {
		node := me.list.prev
		prev := node.prev

		for node != &me.list {
			if err, br := handle(node); br {
				return err, true
			}

			node = prev
			prev = node.prev
		}
	}

	return nil, false
}

/******************************************************************************/

const (
	SizeofHListNode    = 4 * SizeofPointer
	invalidHListBucket = -1
)

type HListNode struct {
	iBucket int
	hash    *HList
	next    *HListNode
	pprev   **HListNode
}

func (me *HListNode) Init() {
	me.iBucket = invalidHListBucket
	me.hash = nil
	me.pprev = nil
	me.next = nil
}

func (me *HListNode) InHash() bool {
	return invalidHListBucket != me.iBucket && nil != me.pprev
}

func (me *HListNode) inTheHash(hash *HList) bool {
	return me.hash == hash
}

func (me *HListNode) InTheHash(hash *HList) bool {
	return me.InHash() && me.inTheHash(hash)
}

func (me *HListNode) del() {
	next := me.next
	pprev := me.pprev

	*pprev = next
	if nil != next {
		next.pprev = pprev
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
	node.Init()
}

type HListForeach func(node *HListNode) (error, bool)

func (me *hlistBucket) foreach(handle HListForeach) (error, bool) {
	if !me.isEmpty() {
		for node := me.first; nil != node; node = node.next {
			if err, br := handle(node); br {
				return err, true
			}
		}
	}

	return nil, false
}

func (me *hlistBucket) unsafeForeach(handle HListForeach) (error, bool) {
	if !me.isEmpty() {
		node := me.first
		next := node.next
		for nil != node {
			if err, br := handle(node); br {
				return err, true
			}

			node = next
			next = node.next
		}
	}

	return nil, false
}

func (me *hlistBucket) find(eq HListEq) *HListNode {
	if !me.isEmpty() {
		for node := me.first; nil != node; node = node.next {
			if eq(node) {
				return node
			}
		}
	}

	return nil
}

func (me *hlistBucket) clean(handle HListForeach) {
	me.unsafeForeach(func(node *HListNode) (error, bool) {
		if nil != handle {
			handle(node)
		}

		me.del(node)

		return nil, false
	})
}

/******************************************************************************/

type HListIndex func() int
type HListEq func(node *HListNode) bool

type HList struct {
	buckets []hlistBucket
	count   int // all node count
}

func (me *HList) Init(count int) {
	me.buckets = make([]hlistBucket, count)
}

func (me *HList) Count() int {
	return me.count
}

func (me *HList) Window() int {
	return len(me.buckets)
}

func (me *HList) bucket(iBucket int) *hlistBucket {
	return &me.buckets[iBucket]
}

func (me *HList) Insert(node *HListNode, index HListIndex) error {
	if node.InHash() {
		Log.Error("insert node:%p into hash:%p error: have in hash:%p", node, me, node.hash)

		return ErrExist
	} else {
		idx := index() % me.Window()
		bucket := me.bucket(idx)
		bucket.add(node)

		node.iBucket = idx
		node.hash = me

		me.count++

		return nil
	}
}

func (me *HList) Remove(node *HListNode, index HListIndex) error {
	if me.count > 0 {
		if node.InHash() {
			if node.inTheHash(me) {

				idx := index() % me.Window()
				bucket := me.bucket(idx)

				if !bucket.isEmpty() {
					bucket.del(node)

					me.count--
				}

				return nil
			} else {
				Log.Error("remove node:%p from hash:%p error: but in the hash:%p",
					node,
					me,
					node.hash)
			}
		} else {
			Log.Error("remove node:%p from hash:%p error: not in the hash", node, me)
		}
	} else {
		Log.Error("remove node:%p from hash:%p error: empty hash", node, me)
	}

	return ErrNoExist
}

func (me *HList) Get(index HListIndex, eq HListEq) *HListNode {
	idx := index() % me.Window()
	bucket := me.bucket(idx)

	return bucket.find(eq)
}

func (me *HList) Clean(handle HListForeach) {
	count := me.Window()

	for i := 0; i < count; i++ {
		bucket := me.bucket(i)

		bucket.clean(handle)
	}

	me.buckets = nil
	me.count = 0
}
