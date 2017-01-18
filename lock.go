package asdf

import (
	"sync"
)

type RwLock struct {
	Debug bool
	Name  string
	m     sync.RWMutex
}

func NewRwLock(name string, debug bool) *RwLock {
	return &RwLock{
		Debug: debug,
		Name:  name,
	}
}

func (me *RwLock) debug(format string, v ...interface{}) {
	if me.Debug {
		Log.Debug(format, v...)
	}
}

func (me *RwLock) rlock() {
	me.debug("%s read lock ...", me.Name)
	me.m.RLock()
	me.debug("%s read lock ok.", me.Name)
}

func (me *RwLock) runlock() {
	me.debug("%s read unlock ...", me.Name)
	me.m.RUnlock()
	me.debug("%s read unlock ok.", me.Name)
}

func (me *RwLock) lock() {
	me.debug("%s write lock ...", me.Name)
	me.m.Lock()
	me.debug("%s write lock ok.", me.Name)
}

func (me *RwLock) unlock() {
	me.debug("%s write unlock ...", me.Name)
	me.m.Unlock()
	me.debug("%s write unlock ok.", me.Name)
}

func (me *RwLock) RHandle(handler func()) {
	me.rlock()
	handler()
	me.runlock()
}

func (me *RwLock) WHandle(handler func()) {
	me.lock()
	handler()
	me.unlock()
}

type Lock struct {
	Debug bool
	Name  string
	m     sync.Mutex
}

func NewLock(name string, debug bool) *Lock {
	return &Lock{
		Debug: debug,
		Name:  name,
	}
}

func (me *Lock) debug(format string, v ...interface{}) {
	if me.Debug {
		Log.Debug(format, v...)
	}
}

func (me *Lock) lock() {
	me.debug("%s lock ...", me.Name)
	me.m.Lock()
	me.debug("%s lock ok.", me.Name)
}

func (me *Lock) unlock() {
	me.debug("%s unlock ...", me.Name)
	me.m.Unlock()
	me.debug("%s unlock ok.", me.Name)
}

func (me *Lock) Handle(handler func()) {
	me.lock()
	handler()
	me.unlock()
}
