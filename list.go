package asdf

type IListFirst interface {
	First() IListNode
}

type IListLast interface {
	Last() IListNode
}

type IListTails interface {
	Tails() IList
}

type IListHeads interface {
	Heads() IList
}

type IListPrev interface {
	Prev() IListNode
}

type IListNext interface {
	Next() IListNode
}

type IListNode interface {
	IListPrev
	IListNext
}

// list = first + tail
// list = head + last
type IList interface {
	IListFirst
	IListLast

	IListTails
	IListHeads

	Insert(IListNode, IListNode, IListNode) error
	Remove(IListNode) error
}

type ISortListNode interface {
	IListNode
	ICompare
}

type ISortList interface {
	IList
}
