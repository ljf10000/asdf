// +build linux

package asdf

import (
	"syscall"
)

func mmap(Len int, inprot, inflags, fd uintptr, off int64) ([]byte, error) {
	flags := syscall.MAP_SHARED
	prot := syscall.PROT_READ

	switch {
	case inprot&COPY != 0:
		prot |= syscall.PROT_WRITE
		flags = syscall.MAP_PRIVATE
	case inprot&RDWR != 0:
		prot |= syscall.PROT_WRITE
	}

	if inprot&EXEC != 0 {
		prot |= syscall.PROT_EXEC
	}

	if inflags&ANON != 0 {
		flags |= syscall.MAP_ANON
	}

	if inflags&HUGE != 0 {
		flags |= syscall.MAP_HUGETLB
	}

	b, err := syscall.Mmap(int(fd), off, Len, prot, flags)
	if err != nil {
		return nil, err
	} else {
		return b, nil
	}
}

func errErrno(errno syscall.Errno) error {
	if 0 == errno {
		return nil
	} else {
		return syscall.Errno(errno)
	}
}

func msync(addr, Len uintptr, sync bool) error {
	flag := syscall.MS_ASYNC
	if sync {
		flag = syscall.MS_SYNC
	}

	_, _, errno := syscall.Syscall(syscall.SYS_MSYNC, addr, Len, uintptr(flag))

	return errErrno(errno)
}

func mlock(addr, Len uintptr) error {
	_, _, errno := syscall.Syscall(syscall.SYS_MLOCK, addr, Len, 0)

	return errErrno(errno)
}

func munlock(addr, Len uintptr) error {
	_, _, errno := syscall.Syscall(syscall.SYS_MUNLOCK, addr, Len, 0)

	return errErrno(errno)
}

func unmap(addr, Len uintptr) error {
	_, _, errno := syscall.Syscall(syscall.SYS_MUNMAP, addr, Len, 0)

	return errErrno(errno)
}
