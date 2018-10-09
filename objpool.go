package asdf

import (
	"fmt"
	"os"
	"path/filepath"
	"unsafe"
)

type objBlock struct {
	file string
	mem  []byte
}

func (me *objBlock) count(size int) int {
	return len(me.mem) / size
}

func (me *objBlock) foreach(size int, handle func(obj *ListNode)) {
	address := SliceAddress(me.mem)

	count := me.count(size)
	for i := 0; i < count; i++ {
		obj := unsafe.Pointer(address + uintptr(i*size))

		handle((*ListNode)(obj))
	}
}

func (me *objBlock) rm() {
	os.Remove(me.file)
}

/******************************************************************************/

type objPoolNode struct {
	ListNode
}

func (me *objPoolNode) Pointer() unsafe.Pointer {
	return unsafe.Pointer(uintptr(unsafe.Pointer(me)) + SizeofListNode)
}

func objPool_obj2node(obj unsafe.Pointer) *objPoolNode {
	return (*objPoolNode)(unsafe.Pointer(uintptr(obj) - SizeofListNode))
}

type ObjPoolOp struct {
	Create func(file string, size int) ([]byte, error)
	Close  func(mem []byte) error
}

type ObjPoolConf struct {
	Dev        string
	Name       string
	ObjSize    int // obj size
	BlockSize  int // block size
	BlockLimit int // block limit
}

type ObjPool struct {
	*ObjPoolOp
	ObjPoolConf

	blockCount int // block count
	blocks     []objBlock

	using List
	free  List
}

func (me *ObjPool) Init(conf *ObjPoolConf, op *ObjPoolOp) error {
	me.ObjPoolOp = op
	me.ObjPoolConf = *conf

	me.mkdir()

	me.blocks = make([]objBlock, 0, me.BlockLimit)

	return nil
}

func (me *ObjPool) Fini() error {
	me.using.Init()
	me.free.Init()

	count := me.blockCount
	for i := 0; i < count; i++ {
		block := me.block(i)

		me.Close(block.mem)
		block.rm()
	}

	me.rmdir()

	return nil
}

func (me *ObjPool) Malloc() (unsafe.Pointer, error) {
	if me.free.IsEmpty() {
		err := me.addPool()
		if nil != err {
			return nil, err
		}
	}

	node := me.free.First()
	me.free.Del(node)
	me.using.Add(node)

	obj := (*objPoolNode)(unsafe.Pointer(node))

	return obj.Pointer(), nil
}

func (me *ObjPool) Free(obj unsafe.Pointer) {
	node := objPool_obj2node(obj)

	me.using.Del(&node.ListNode)
	me.free.Add(&node.ListNode)
}

func (me *ObjPool) dir() string {
	return filepath.Join(me.Dev, me.Name)
}

func (me *ObjPool) mkdir() {
	FileName(me.dir()).Mkdir()
}

func (me *ObjPool) rmdir() {
	os.RemoveAll(me.dir())
}

func (me *ObjPool) file(iBlock int) string {
	file := fmt.Sprintf("04x", iBlock)

	return filepath.Join(me.dir(), file)
}

func (me *ObjPool) objSize() int {
	return me.ObjSize + SizeofListNode
}

func (me *ObjPool) block(iBlock int) *objBlock {
	return &me.blocks[iBlock]
}

func (me *ObjPool) newBlock(iBlock int) error {
	file := me.file(iBlock)

	mem, err := me.Create(file, me.BlockSize)
	if nil != err {
		return err
	}

	block := me.block(iBlock)
	block.file = file
	block.mem = mem

	return nil
}

func (me *ObjPool) addPool() error {
	iBlock := me.blockCount
	if iBlock == me.BlockLimit {
		return ErrNoSpace
	}

	err := me.newBlock(iBlock)
	if nil != err {
		return err
	}
	me.blockCount++

	block := me.block(iBlock)
	block.foreach(me.objSize(), func(obj *ListNode) {
		me.free.Add(obj)
	})

	return nil
}
