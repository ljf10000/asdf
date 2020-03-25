package asdf

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"
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
	Lock       bool
	Dev        string
	Name       string
	ObjSize    int // obj size
	BlockSize  int // block size
	BlockLimit int // block limit
}

func (me *ObjPoolConf) file(iBlock int) string {
	file := fmt.Sprintf("%s-%04x", me.Name, iBlock)

	return filepath.Join(me.Dev, file)
}

/******************************************************************************/

type objPoolList struct {
	times uint64
	list  List
}

type ObjPoolStat struct {
	Dev  string `json:"dev"`
	Name string `json:"name"`

	Block struct {
		Limit int `json:"limit"`
		Count int `json:"count"`
		Size  int `json:"size"`
	} `json:"block"`

	Obj struct {
		Size  int `json:"size"`
		Using struct {
			Times int `json:"times"`
			Count int `json:"count"`
		} `json:"using"`
		Freed struct {
			Times int `json:"times"`
			Count int `json:"count"`
		} `json:"freed"`
	} `json:"obj"`
}

func (me *ObjPoolStat) String() string {
	return fmt.Sprintf("using[%d:%d] freed[%d:%d]",
		me.Obj.Using.Times,
		me.Obj.Using.Count,
		me.Obj.Freed.Times,
		me.Obj.Freed.Count)
}

func (me *ObjPool) Stat() ObjPoolStat {
	stat := ObjPoolStat{
		Name: me.Name,
		Dev:  me.Dev,
	}

	stat.Block.Limit = me.BlockLimit
	stat.Block.Size = me.BlockSize

	stat.Obj.Size = me.ObjSize

	return stat
}

func (me *ObjPool) UpdateStat(stat *ObjPoolStat) {
	stat.Block.Count = me.blockCount

	stat.Obj.Using.Times = int(me.using.times)
	stat.Obj.Using.Count = me.using.list.count
	stat.Obj.Freed.Times = int(me.freed.times)
	stat.Obj.Freed.Count = me.freed.list.count
}

func ObjPoolPrepare(root, prefix string) error {
	Prefix := filepath.Join(root, prefix)

	return filepath.Walk(root, func(path string, f os.FileInfo, err error) error {
		if nil == f {
			return err
		}

		if !f.IsDir() && strings.HasPrefix(path, Prefix) {
			FileName(path).Delete()
		}

		return nil
	})
}

type ObjPool struct {
	*ObjPoolOps
	*ObjPoolConf

	m sync.Mutex

	blockCount int // block count
	blocks     []ObjPoolBlock

	using objPoolList
	freed objPoolList
}

func (me *ObjPool) lock() {
	if me.Lock {
		me.m.Lock()
	}
}

func (me *ObjPool) unlock() {
	if me.Lock {
		me.m.Unlock()
	}
}

func (me *ObjPool) handle(cb func()) {
	me.lock()
	cb()
	me.unlock()
}

func (me *ObjPool) Init(conf *ObjPoolConf, ops *ObjPoolOps) error {
	me.ObjPoolOps = ops
	me.ObjPoolConf = conf

	me.handle(func() {
		me.using.list.Init()
		me.freed.list.Init()
		me.blocks = make([]ObjPoolBlock, me.BlockLimit)
	})

	Log.Debug("objpool %s init", me.Name)

	return nil
}

func (me *ObjPool) Fini() error {
	me.handle(func() {
		count := me.blockCount
		for i := 0; i < count; i++ {
			block := me.block(i)

			block.fini(me.ObjPoolOps)
		}

		me.using.list.Init()
		me.freed.list.Init()
	})

	Log.Debug("objpool %s fini", me.Name)

	return nil
}

func (me *ObjPool) Malloc() (ptr unsafe.Pointer, err error) {
	me.handle(func() {
		err = me.preMalloc()
		if nil == err {
			obj := me.malloc()

			ptr = GetObjField(obj, SizeofListNode)
		}
	})

	return ptr, err
}

func (me *ObjPool) Free(obj unsafe.Pointer) {
	me.handle(func() {
		node := GetObjByField(obj, SizeofListNode)

		me.free(node)
	})
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
