// +build linux

package asdf

import (
	"os"
	"syscall"
)

func (me FileName) Lock() error {
	f, err := os.OpenFile(string(me), os.O_RDWR|os.O_CREATE, 0666)
	if nil != err {
		return err
	}

	err = syscall.Flock(int(f.Fd()), syscall.LOCK_EX|syscall.LOCK_NB)
	if nil != err {
		return err
	}

	return nil
}
