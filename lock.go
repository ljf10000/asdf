package asdf

import (
	"sync"
)

type ILock interface {
	Write(handler func() error) error
	Read(handler func() error) error
}

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

func (me *RwLock) Read(handler func() error) error {
	var err error

	me.rlock()
	err = handler()
	me.runlock()

	return err
}

func (me *RwLock) Write(handler func() error) error {
	var err error

	me.lock()
	err = handler()
	me.unlock()

	return err
}

type AccessLock struct {
	Debug bool
	Name  string
	m     sync.Mutex
}

func NewAccessLock(name string, debug bool) *AccessLock {
	return &AccessLock{
		Debug: debug,
		Name:  name,
	}
}

func (me *AccessLock) debug(format string, v ...interface{}) {
	if me.Debug {
		Log.Debug(format, v...)
	}
}

func (me *AccessLock) lock() {
	me.debug("%s lock ...", me.Name)
	me.m.Lock()
	me.debug("%s lock ok.", me.Name)
}

func (me *AccessLock) unlock() {
	me.debug("%s unlock ...", me.Name)
	me.m.Unlock()
	me.debug("%s unlock ok.", me.Name)
}

func (me *AccessLock) Handle(handler func()) {
	me.lock()
	handler()
	me.unlock()
}

func (me *AccessLock) RHandle(handler func()) {
	me.Handle(handler)
}

func (me *AccessLock) WHandle(handler func()) {
	me.Handle(handler)
}

func (me *AccessLock) HandleE(handler func() error) error {
	var err error

	me.lock()
	err = handler()
	me.unlock()

	return err
}

func (me *AccessLock) Read(handler func() error) error {
	return me.HandleE(handler)
}

func (me *AccessLock) Write(handler func() error) error {
	return me.HandleE(handler)
}
