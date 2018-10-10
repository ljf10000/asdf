package asdf

import (
	"unsafe"
)

/******************************************************************************/

const (
	SizeofHashNode   = 4 * SizeofPointer
	invalidHashBucket = -1
)

var scHashNode = NewSizeChecker("HashNode", unsafe.Sizeof(HashNode{}), SizeofHashNode)

type HashNode struct {
	iBucket int
	hash    *Hash
	next    *HashNode
	pprev   **HashNode
}

func (me *HashNode) Init() {
	me.iBucket = invalidHashBucket
	me.hash = nil
	me.pprev = nil
	me.next = nil
}

func (me *HashNode) InHash() bool {
	return invalidHashBucket != me.iBucket && nil != me.pprev
}

func (me *HashNode) inTheHash(hash *Hash) bool {
	return me.hash == hash
}

func (me *HashNode) InTheHash(hash *Hash) bool {
	return me.InHash() && me.inTheHash(hash)
}

func (me *HashNode) del() {
	next := me.next
	pprev := me.pprev

	*pprev = next
	if nil != next {
		next.pprev = pprev
	}
}

/******************************************************************************/

type hashBucket struct {
	first *HashNode
}

func (me *hashBucket) isEmpty() bool {
	return nil == me.first
}

func (me *hashBucket) add(node *HashNode) {
	first := me.first

	node.next = first
	if nil != first {
		first.pprev = &node.next
	}
	me.first = node
	node.pprev = &me.first
}

func (me *hashBucket) del(node *HashNode) {
	node.del()
	node.Init()
}

type HListForeach func(node *HashNode) (error, bool)

func (me *hashBucket) foreach(handle HListForeach) (error, bool) {
	if !me.isEmpty() {
		for node := me.first; nil != node; node = node.next {
			if err, br := handle(node); br {
				return err, true
			}
		}
	}

	return nil, false
}

func (me *hashBucket) unsafeForeach(handle HListForeach) (error, bool) {
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

func (me *hashBucket) find(eq HListEq) *HashNode {
	if !me.isEmpty() {
		for node := me.first; nil != node; node = node.next {
			if eq(node) {
				return node
			}
		}
	}

	return nil
}

func (me *hashBucket) clean(handle HListForeach) {
	me.unsafeForeach(func(node *HashNode) (error, bool) {
		if nil != handle {
			handle(node)
		}

		me.del(node)

		return nil, false
	})
}

/******************************************************************************/

type HListIndex func() int
type HListEq func(node *HashNode) bool

type Hash struct {
	buckets []hashBucket
	count   int // all node count
}

func (me *Hash) Init(count int) {
	me.buckets = make([]hashBucket, count)
}

func (me *Hash) Count() int {
	return me.count
}

func (me *Hash) Window() int {
	return len(me.buckets)
}

func (me *Hash) bucket(iBucket int) *hashBucket {
	return &me.buckets[iBucket]
}

func (me *Hash) Insert(node *HashNode, index HListIndex) error {
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

func (me *Hash) Remove(node *HashNode, index HListIndex) error {
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

func (me *Hash) Get(index HListIndex, eq HListEq) *HashNode {
	idx := index() % me.Window()
	bucket := me.bucket(idx)

	return bucket.find(eq)
}

func (me *Hash) Clean(handle HListForeach) {
	count := me.Window()

	for i := 0; i < count; i++ {
		bucket := me.bucket(i)

		bucket.clean(handle)
	}

	me.buckets = nil
	me.count = 0
}
