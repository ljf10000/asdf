// +build linux

package asdf

import (
	"os"
	"syscall"
)

func fileLock(file string) error {
	// readonly
	f, err := os.OpenFile(file, os.O_RDONLY, 0666)
	if nil != err {
		return err
	}

	err = syscall.Flock(int(f.Fd()), syscall.LOCK_EX|syscall.LOCK_NB)
	if nil != err {
		return err
	}

	return nil
}
