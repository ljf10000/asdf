package asdf

import (
	"sync/atomic"
	"unsafe"
)

type MpLockRw struct {
	r mplock
	w mplock
}

func (me *MpLockRw) dump(role mplock_role_t, act mplock_act_t, tag string) {
	if MpLockDump {
		Log.Info("%s role[%s] act[%s]: reader[want:%d, value:%d] writer[want:%d, value:%d]",
			tag,
			role,
			act,
			me.r.want(), me.r.value(),
			me.w.want(), me.w.value())
	}
}

func (me *MpLockRw) split(role mplock_role_t, self *mplock, other *mplock) {
	if MPLOCK_READER == role {
		self.v = me.r.v
		other.v = me.w.v
	} else {
		self.v = me.w.v
		other.v = me.r.v
	}
}

func (me *MpLockRw) merge(role mplock_role_t, self *mplock, other *mplock) {
	if MPLOCK_READER == role {
		me.r.v = self.v
		me.w.v = other.v
	} else {
		me.w.v = self.v
		me.r.v = other.v
	}
}

type mplock_rw_try_f func() bool
type mplock_rw_do_f func(self *mplock, other *mplock)

func (me *MpLockRw) until_ok(try mplock_rw_try_f) bool {
	for {
		if try() {
			return true
		} else {
			MpLockPause()
		}
	}

	return false
}

func (me *MpLockRw) try(
	role mplock_role_t,
	act mplock_act_t,
	can mplock_can_f,
	exec mplock_rw_do_f,
	want mplock_rw_do_f) bool {
	//--------------------------------------------------------------------------
	var new_, old_ MpLockRw
	var self, other mplock

	v := atomic.LoadUint32(me.Address())

	old_.Update(v)
	old_.split(role, &self, &other)

	if err := mplock_check(role, &self, &other); nil != err {
		return false
	}

	switch can(role, &self, &other) {
	case MPLOCK_CAN_EXEC:
		exec(&self, &other)

		new_.merge(role, &self, &other)
		if atomic.CompareAndSwapUint32(me.Address(), v, new_.Value()) {
			old_.dump(role, act, "mplock-rw exec old")
			new_.dump(role, act, "mplock-rw exec new")

			return true
		}
	case MPLOCK_CAN_WANT:
		want(&self, &other)

		new_.merge(role, &self, &other)
		if atomic.CompareAndSwapUint32(me.Address(), v, new_.Value()) {
			old_.dump(role, act, "mplock-rw want old")
			new_.dump(role, act, "mplock-rw want new")

			// do want, NOT do lock, return false
		}
	}

	return false
}

func mplock_rw_want(self *mplock, other *mplock) {
	if other.want() {
		Panic("other have wanted")
	}

	self.set_want(true)
}

func mplock_rw_lock(self *mplock, other *mplock) {
	if other.want() {
		Panic("other have wanted")
	} else if other.locked() {
		Panic("other have locked")
	}

	self.add()
	self.set_want(false)
}

func mplock_rw_unlock(self *mplock, other *mplock) {
	self.sub()
}

func mplock_rw_upgrade(self *mplock, other *mplock) {
	self.add()
	other.sub()
}

func mplock_rw_degrade(self *mplock, other *mplock) {
	self.add()
	other.sub()
}

func (me *MpLockRw) lock(role mplock_role_t) bool {
	return me.try(role, MPLOCK_ACT_LOCK, mplock_can_lock, mplock_rw_lock, mplock_rw_want)
}

func (me *MpLockRw) unlock(role mplock_role_t) bool {
	return me.try(role, MPLOCK_ACT_UNLOCK, mplock_can_unlock, mplock_rw_unlock, nil)
}

func (me *MpLockRw) upgrade(role mplock_role_t) bool {
	return me.try(role, MPLOCK_ACT_UPGRADE, mplock_can_upgrade, mplock_rw_upgrade, nil)
}

func (me *MpLockRw) degrade(role mplock_role_t) bool {
	return me.try(role, MPLOCK_ACT_DEGRADE, mplock_can_degrade, mplock_rw_degrade, nil)
}

func (me *MpLockRw) Address() *uint32 {
	return (*uint32)(unsafe.Pointer(me))
}

func (me *MpLockRw) Value() uint32 {
	return *me.Address()
}

func (me *MpLockRw) Update(v uint32) {
	*me.Address() = v
}

func (me *MpLockRw) Init() {
	me.Update(0)
}

func (me *MpLockRw) LockR() bool {
	return me.until_ok(func() bool {
		return me.lock(MPLOCK_READER)
	})
}

func (me *MpLockRw) LockW() bool {
	return me.until_ok(func() bool {
		return me.lock(MPLOCK_WRITER)
	})
}

func (me *MpLockRw) UnLockR() bool {
	return me.until_ok(func() bool {
		return me.unlock(MPLOCK_READER)
	})
}

func (me *MpLockRw) UnLockW() bool {
	return me.until_ok(func() bool {
		return me.unlock(MPLOCK_WRITER)
	})
}

func (me *MpLockRw) Upgrade_rtow() bool {
	return me.until_ok(func() bool {
		return me.upgrade(MPLOCK_READER)
	})
}

func (me *MpLockRw) Upgrade_wtor() bool {
	return me.until_ok(func() bool {
		return me.upgrade(MPLOCK_WRITER)
	})
}

func (me *MpLockRw) Degrade_rtow() bool {
	return me.until_ok(func() bool {
		return me.degrade(MPLOCK_READER)
	})
}

func (me *MpLockRw) Degrade_wtor() bool {
	return me.until_ok(func() bool {
		return me.degrade(MPLOCK_WRITER)
	})
}

func (me *MpLockRw) HandleR(handle func()) {
	if me.LockR() {
		handle()

		me.UnLockR()
	}
}

func (me *MpLockRw) HandleW(handle func()) {
	if me.LockW() {
		handle()

		me.UnLockW()
	}
}
