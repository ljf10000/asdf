package asdf

// changed from https://github.com/edsrzf/mmap-go
//	liujingfei
// 	support msync(async)
//	not save windows HANDLE(file)

import (
	"errors"
	"os"
	"reflect"
	"unsafe"
)

/******************************************************************************/

const (
	MMAP_RO   = 0
	MMAP_RW   = 0x01
	MMAP_COPY = 0x02
	MMAP_EXEC = 0x04
)

const (
	MMAP_ANON = 1 << iota
	MMAP_HUGE = 1 << iota
)

// MMap represents a file mapped into memory.
type MMap []byte

const (
	MMAP_F_CREATE = 0x01
	MMAP_F_ADJUST = 0x02
	MMAP_F_HUGE   = 0x04
	MMAP_F_RO     = 0x08
)

// flag is ro/huge
func MapOpenEx(file string, size int, flag int) (MMap, error) {
	if FileName(file).Exist() {
		return MapOpen(file, flag)
	} else {
		return MapCreate(file, size, flag)
	}
}

// flag is ro/huge
func MapCreate(file string, size int, flag int) (MMap, error) {
	return MapFile(file, size, flag|MMAP_F_CREATE)
}

// flag is ro/huge
func MapOpen(file string, flag int) (MMap, error) {
	return MapFile(file, -1, flag)
}

// flag is ro/huge
func MapTrunc(file string, size int, flag int) (MMap, error) {
	return MapFile(file, size, flag|MMAP_F_ADJUST)
}

func MapFile(file string, size int, flag int) (MMap, error) {
	readonly := MMAP_F_RO == (flag & MMAP_F_RO)
	create := MMAP_F_CREATE == (flag & MMAP_F_CREATE)
	adj := MMAP_F_ADJUST == (flag & MMAP_F_ADJUST)

	openflag := 0
	prot := 0
	if !readonly {
		openflag = os.O_RDWR
		prot = MMAP_RW
	}

	if size < 0 {
		if adj {
			return nil, ErrSprintf("mapfile adj with size<0")
		}

		if create {
			return nil, ErrSprintf("mapfile create with size<0")
		}
	}

	if create {
		openflag |= os.O_CREATE | os.O_TRUNC
	}

	f, err := os.OpenFile(file, openflag, 0666)
	if nil != err {
		Log.Error("open %s error:%s", file, err)

		return nil, err
	}

	//Log.Debug("OK: open %s", file)

	defer func() {
		f.Close()
		//Log.Debug("OK: close %s", file)
	}()

	if size < 0 {
		size = int(FileSize(f))
	} else {
		size = PageAlign(size)
	}

	// 1. create new file
	// 2. open exist file, and adjust
	if create || adj {
		err := f.Truncate(int64(size))
		if nil != err {
			Log.Error("truncate %s size:%d error:%s", file, size, err)

			return nil, err
		}

		Log.Debug("OK: truncate %s size:%d", file, size)
	}

	//Log.Debug("OK: mmap %s size:%d", file, size)

	// try use hugepage
	inflags := 0
	if MMAP_F_HUGE == (flag & MMAP_F_HUGE) {
		// inflags = HUGE
	}

	return MapFileEx(f, file, size, prot, inflags, 0)
}

func MapFileEx(f *os.File, file string, size, prot, flag int, offset int64) (MMap, error) {
	if offset%int64(PAGESIZE) != 0 {
		return nil, errors.New("offset parameter must be a multiple of the system's page size")
	}

	var fd uintptr

	if 0 == MMAP_ANON&flag {
		fd = uintptr(f.Fd())
	} else {
		if size <= 0 {
			return nil, errors.New("anonymous mapping requires non-zero length")
		}

		fd = ^uintptr(0)
	}

	return mmap(size, uintptr(prot), uintptr(flag), fd, offset)
}

func (me *MMap) header() *reflect.SliceHeader {
	return (*reflect.SliceHeader)(unsafe.Pointer(me))
}

func (me MMap) Zero() {
	count := len(me)
	for i := 0; i < count; i++ {
		me[i] = 0
	}
}

func (me MMap) Lock() error {
	h := me.header()

	return mlock(h.Data, uintptr(h.Len))
}

func (me MMap) Unlock() error {
	h := me.header()

	return munlock(h.Data, uintptr(h.Len))
}

func (me MMap) Msync(sync bool) error {
	h := me.header()

	return msync(h.Data, uintptr(h.Len), sync)
}

func (me MMap) Unmap() error {
	h := me.header()

	err := unmap(h.Data, uintptr(h.Len))
	if nil != err {
		return err
	}

	return nil
}
