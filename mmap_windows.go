// +build windows

package asdf

import (
	"os"
	"syscall"
)

func mmap(Len int, prot, flags, hfile uintptr, off int64) ([]byte, error) {
	flProtect := uint32(syscall.PAGE_READONLY)
	dwDesiredAccess := uint32(syscall.FILE_MAP_READ)
	switch {
	case prot&MMAP_COPY != 0:
		flProtect = syscall.PAGE_WRITECOPY
		dwDesiredAccess = syscall.FILE_MAP_COPY
	case prot&MMAP_RW != 0:
		flProtect = syscall.PAGE_READWRITE
		dwDesiredAccess = syscall.FILE_MAP_WRITE
	}
	if prot&MMAP_EXEC != 0 {
		flProtect <<= 4
		dwDesiredAccess |= syscall.FILE_MAP_EXECUTE
	}

	// The maximum size is the area of the file, starting from 0,
	// that we wish to allow to be mappable. It is the sum of
	// the length the user requested, plus the offset where that length
	// is starting from. This does not map the data into memory.
	maxSizeHigh := uint32((off + int64(Len)) >> 32)
	maxSizeLow := uint32((off + int64(Len)) & 0xFFFFFFFF)

	hMap, errno := syscall.CreateFileMapping(syscall.Handle(hfile), nil, flProtect, maxSizeHigh, maxSizeLow, nil)
	if hMap == 0 {
		return nil, os.NewSyscallError("CreateFileMapping", errno)
	}
	defer syscall.CloseHandle(hMap)

	// Actually map a view of the data into memory. The view's size
	// is the length the user requested.
	fileOffsetHigh := uint32(off >> 32)
	fileOffsetLow := uint32(off & 0xFFFFFFFF)
	addr, errno := syscall.MapViewOfFile(hMap, dwDesiredAccess, fileOffsetHigh, fileOffsetLow, uintptr(Len))
	if addr == 0 {
		return nil, os.NewSyscallError("MapViewOfFile", errno)
	}

	m := MMap{}
	h := m.header()
	h.Data = addr
	h.Len = Len
	h.Cap = h.Len

	return m, nil
}

// windows not support async
func msync(addr, Len uintptr, sync bool) error {
	errno := syscall.FlushViewOfFile(addr, Len)

	return os.NewSyscallError("FlushViewOfFile", errno)
}

func mlock(addr, Len uintptr) error {
	errno := syscall.VirtualLock(addr, Len)

	return os.NewSyscallError("VirtualLock", errno)
}

func munlock(addr, Len uintptr) error {
	errno := syscall.VirtualUnlock(addr, Len)

	return os.NewSyscallError("VirtualUnlock", errno)
}

func unmap(addr, Len uintptr) error {
	msync(addr, Len, true)

	errno := syscall.UnmapViewOfFile(addr)

	return os.NewSyscallError("UnmapViewOfFile", errno)
}
