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
