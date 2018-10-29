package asdf

import (
	"fmt"
	"os"
	"path/filepath"
	"unsafe"
)

func GetObjField(obj unsafe.Pointer, offset uintptr) unsafe.Pointer {
	return unsafe.Pointer(uintptr(obj) + offset)
}

func GetObjByField(field unsafe.Pointer, offset uintptr) unsafe.Pointer {
	return unsafe.Pointer(uintptr(field) - offset)
}

/******************************************************************************/

type ObjPoolBlock struct {
	file string
	mem  []byte
}

func (me *ObjPoolBlock) count(objSize int) int {
	return len(me.mem) / objSize
}

func (me *ObjPoolBlock) foreach(objSize int, handle func(obj *ListNode)) {
	address := SliceAddress(me.mem)

	count := me.count(objSize)
	for i := 0; i < count; i++ {
		obj := (*ListNode)(unsafe.Pointer(address + uintptr(i*objSize)))

		handle(obj)
	}
}

func (me *ObjPoolBlock) fini(ops *ObjPoolOps) {
	ops.Close(me.mem)
	me.mem = nil

	os.Remove(me.file)
}

/******************************************************************************/

type objPoolNode = ListNode

type ObjPoolOps struct {
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

func (me *ObjPoolConf) file(iBlock int) string {
	dir := me.dir()
	file := fmt.Sprintf("04x", iBlock)

	return filepath.Join(dir, file)
}

func (me *ObjPoolConf) dir() string {
	return filepath.Join(me.Dev, me.Name)
}

func (me *ObjPoolConf) mkdir() error {
	dir := me.dir()

	err := FileName(dir).Mkdir()
	if nil != err {
		Log.Error("obj pool(%s) mkdir(%s) error:%s", me.Name, dir, err)

		return err
	}

	return nil
}

func (me *ObjPoolConf) rmdir() error {
	dir := me.dir()

	err := os.RemoveAll(dir)
	if nil != err {
		Log.Error("obj pool(%s) rmdir(%s) error:%s", me.Name, dir, err)

		return err
	}

	return nil
}

/******************************************************************************/

type objPoolList struct {
	times uint64
	list  List
}

type ObjPool struct {
	*ObjPoolOps
	*ObjPoolConf

	blockCount int // block count
	blocks     []ObjPoolBlock

	using objPoolList
	freed objPoolList
}

func (me *ObjPool) Init(conf *ObjPoolConf, ops *ObjPoolOps) error {
	me.ObjPoolOps = ops
	me.ObjPoolConf = conf

	me.blocks = make([]ObjPoolBlock, me.BlockLimit)

	return me.mkdir()
}

func (me *ObjPool) Fini() error {
	me.using.list.Init()
	me.freed.list.Init()

	count := me.blockCount
	for i := 0; i < count; i++ {
		block := me.block(i)

		block.fini(me.ObjPoolOps)
	}

	return me.rmdir()
}

func (me *ObjPool) Malloc() (unsafe.Pointer, error) {
	err := me.preMalloc()
	if nil != err {
		return nil, err
	}

	obj := me.malloc()

	return GetObjField(obj, SizeofListNode), nil
}

func (me *ObjPool) Free(obj unsafe.Pointer) {
	node := GetObjByField(obj, SizeofListNode)

	me.free(node)
}

/******************************************************************************/

func (me *ObjPool) preMalloc() error {
	if me.freed.list.IsEmpty() {
		err := me.addPool()
		if nil != err {
			return err
		}
	}

	return nil
}

func (me *ObjPool) malloc() unsafe.Pointer {
	node := me.freed.list.First()

	me.freed.list.Remove(node)
	me.using.list.InsertHead(node)
	me.using.times++

	return unsafe.Pointer(node)
}

func (me *ObjPool) free(obj unsafe.Pointer) {
	node := (*objPoolNode)(obj)

	me.using.list.Remove(node)
	me.freed.list.InsertHead(node)
	me.freed.times++
}

func (me *ObjPool) objSize() int {
	return me.ObjSize + SizeofListNode
}

func (me *ObjPool) block(iBlock int) *ObjPoolBlock {
	return &me.blocks[iBlock]
}

func (me *ObjPool) newBlock() (*ObjPoolBlock, error) {
	iBlock := me.blockCount
	file := me.file(iBlock)

	mem, err := me.Create(file, me.BlockSize)
	if nil != err {
		return nil, err
	}

	block := me.block(iBlock)
	block.file = file
	block.mem = mem

	me.blockCount++

	return block, nil
}

func (me *ObjPool) addPool() error {
	if me.isFullBlock() {
		return ErrLimit
	}

	block, err := me.newBlock()
	if nil != err {
		return err
	}

	me.addBlock(block)

	return nil
}

func (me *ObjPool) addBlock(block *ObjPoolBlock) {
	block.foreach(me.objSize(), func(obj *ListNode) {
		me.freed.list.InsertHead(obj)
	})
}

func (me *ObjPool) isFullBlock() bool {
	return me.blockCount == me.BlockLimit
}
