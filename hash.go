package asdf

import (
	"container/list"
)

type IHashData interface {
	HashCode() int
	CreateNode() (IHashNode, error)
}

type IHashNode interface {
	HashCode() int
	Eq(data IHashData) bool
	Fini()

	GetElm() *list.Element
	SetElm(elm *list.Element)

	SetHashCode(code int)
	GetHashCode() int
}

type Hash struct {
	buckets []list.List
}

func NewHash(count int) *Hash {
	obj := &Hash{
		buckets: make([]list.List, count),
	}

	for i := 0; i < count; i++ {
		obj.buckets[i].Init()
	}

	return obj
}

func (me *Hash) bucket(code int) *list.List {
	return &me.buckets[code%len(me.buckets)]
}

func (me *Hash) find(code int, data IHashData) IHashNode {
	lst := me.bucket(code)

	for elm := lst.Front(); nil != elm; elm = elm.Next() {
		node, _ := elm.Value.(IHashNode)

		if node.Eq(data) {
			return node
		}
	}

	return nil
}

func (me *Hash) Find(data IHashData) IHashNode {
	return me.find(data.HashCode(), data)
}

func (me *Hash) Remove(node IHashNode) {
	elm := node.GetElm()
	if nil != elm {
		lst := me.bucket(node.GetHashCode())

		lst.Remove(elm)
		node.SetElm(nil)
		node.SetHashCode(0)
	}
}

func (me *Hash) Insert(data IHashData) IHashNode {
	code := data.HashCode()

	if node := me.find(code, data); nil != node {
		return node
	}

	node, err := data.CreateNode()
	if nil != err {
		return nil
	}

	lst := me.bucket(code)
	elm := lst.PushFront(node)
	node.SetElm(elm)
	node.SetHashCode(code)

	return node
}

func (me *Hash) Clean() {
	count := len(me.buckets)

	for i := 0; i < count; i++ {
		lst := me.bucket(i)
		elm := lst.Front()

		for nil != elm {
			v := lst.Remove(elm)

			node, _ := v.(IHashNode)

			node.Fini()
			node.SetElm(nil)
			node.SetHashCode(0)

			elm = lst.Front()
		}
	}
}

func (me *Hash) Fini() {
	me.buckets = nil
}
