// +build linux

package asdf

import (
	"syscall"
)

func (me FileName) Lock() error {
	f, err := os.OpenFile(string(me), os.O_RDWR, 0666)
	if nil != err {
		Log.Error("open %s error: %s", me, err)

		return err
	}

	err = syscall.Flock(f.Fd(), syscall.LOCK_EX|syscall.LOCK_NB)
	if nil != err {
		Log.Error("lock %s error: %s", me, err)

		return err
	}

	return nil
}
