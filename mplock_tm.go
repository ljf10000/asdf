package asdf

import (
	"sync/atomic"
	"unsafe"
)

type MpLockTm struct {
	time Time32

	MpLockRw
}

func (me *MpLockTm) dump(role mplock_role_t, act mplock_act_t, tag string) {
	if MpLockDump {
		Log.Info("%s role[%s] act[%s]: time[%u] reader[want:%d, value:%d] writer[want:%d, value:%d]",
			tag,
			role,
			act,
			me.time,
			me.r.want(), me.r.value(),
			me.w.want(), me.w.value())
	}
}

type mplock_tm_try_f func() MpLockState
type mplock_tm_do_f func(lock *MpLockTm, self *mplock, other *mplock)

func (me *MpLockTm) until_ok(try mplock_tm_try_f) MpLockState {
	for {
		state := try()
		if MPLOCK_SECCESS == state || MPLOCK_TIMEOUT == state {
			return state
		} else {
			MpLockPause()
		}
	}

	return MPLOCK_FAILED
}

func (me *MpLockTm) need_check_timeout(role mplock_role_t) bool {
	return (MPLOCK_READER == role && me.w.locked()) ||
		(MPLOCK_WRITER == role && me.r.locked())
}

func (me *MpLockTm) is_locked() bool {
	return me.r.locked() || me.w.locked()
}

func (me *MpLockTm) lifetime() Time32 {
	return mplock_time() - me.time
}

func (me *MpLockTm) is_timeout(timeout int) bool {
	return timeout > 0 &&
		me.time > 0 &&
		me.lifetime() > Time32(timeout)
}

func (me *MpLockTm) try(
	role mplock_role_t,
	act mplock_act_t,
	timeout int,
	can mplock_can_f,
	exec mplock_tm_do_f,
	want mplock_tm_do_f) MpLockState {
	//--------------------------------------------------------------------------
	var new_, old_ MpLockTm
	var self, other mplock

	v := atomic.LoadUint64(me.Address())
	old_.Update(v)
	old_.split(role, &self, &other)

	// Log.Debug("mplock tm %s try %s:\n%s", role, act, BinSprintf(StructSlice(unsafe.Pointer(me), 8)))

	if err := mplock_check(role, &self, &other); nil != err {
		return MPLOCK_FAILED
	}

	if old_.need_check_timeout(role) && old_.is_timeout(timeout) {
		return MPLOCK_TIMEOUT
	}

	switch can(role, &self, &other) {
	case MPLOCK_CAN_EXEC:
		exec(&new_, &self, &other)

		new_.merge(role, &self, &other)

		if atomic.CompareAndSwapUint64(me.Address(), v, new_.Value()) {
			old_.dump(role, act, "mplock-tm exec old")
			new_.dump(role, act, "mplock-tm exec new")

			return MPLOCK_SECCESS
		}

	case MPLOCK_CAN_WANT:
		want(&new_, &self, &other)

		new_.merge(role, &self, &other)
		if atomic.CompareAndSwapUint64(me.Address(), v, new_.Value()) {
			old_.dump(role, act, "mplock-tm want old")
			new_.dump(role, act, "mplock-tm want new")

			// do want, NOT do lock, return false
		}
	}

	return MPLOCK_FAILED
}

func mplock_tm_want(lock *MpLockTm, self *mplock, other *mplock) {
	mplock_rw_want(self, other)
}

func mplock_tm_lock(lock *MpLockTm, self *mplock, other *mplock) {
	if self.not_locked() {
		lock.time = mplock_time()
	}

	mplock_rw_lock(self, other)
}

func mplock_tm_unlock(lock *MpLockTm, self *mplock, other *mplock) {
	mplock_rw_unlock(self, other)

	if self.not_locked() {
		lock.time = 0
	}
}

func mplock_tm_upgrade(lock *MpLockTm, self *mplock, other *mplock) {
	if self.not_locked() {
		lock.time = mplock_time()
	}

	mplock_rw_upgrade(self, other)
}

func mplock_tm_degrade(lock *MpLockTm, self *mplock, other *mplock) {
	mplock_rw_degrade(self, other)

	if self.not_locked() {
		lock.time = 0
	}
}

func (me *MpLockTm) lock(role mplock_role_t, timeout int) MpLockState {
	return me.try(role, MPLOCK_ACT_LOCK, timeout, mplock_can_lock, mplock_tm_lock, mplock_tm_want)
}

func (me *MpLockTm) unlock(role mplock_role_t) MpLockState {
	return me.try(role, MPLOCK_ACT_UNLOCK, 0, mplock_can_unlock, mplock_tm_unlock, nil)
}

func (me *MpLockTm) upgrade(role mplock_role_t, timeout int) MpLockState {
	return me.try(role, MPLOCK_ACT_UPGRADE, timeout, mplock_can_upgrade, mplock_tm_upgrade, nil)
}

func (me *MpLockTm) degrade(role mplock_role_t, timeout int) MpLockState {
	return me.try(role, MPLOCK_ACT_DEGRADE, timeout, mplock_can_degrade, mplock_tm_degrade, nil)
}

func (me *MpLockTm) recover_() bool {
	var tmp MpLockTm

	v := atomic.LoadUint64(me.Address())
	tmp.Update(v)

	tmp.time = 0
	if tmp.r.locked() {
		tmp.r.sub()
		tmp.r.set_want(false)
	} else if tmp.w.locked() {
		tmp.w.sub()
		tmp.w.set_want(false)
	}

	return atomic.CompareAndSwapUint64(me.Address(), v, tmp.Value())
}

func (me *MpLockTm) Recover(timeout int) bool {
	var tmp MpLockTm

	v := atomic.LoadUint64(me.Address())
	tmp.Update(v)

	is_timeout := tmp.is_locked() && tmp.is_timeout(timeout)
	if is_timeout {
		me.until_ok(func() MpLockState {
			if me.recover_() {
				return MPLOCK_SECCESS
			} else {
				return MPLOCK_FAILED
			}
		})
	}

	return is_timeout
}

func (me *MpLockTm) Address() *uint64 {
	return (*uint64)(unsafe.Pointer(me))
}

func (me *MpLockTm) Value() uint64 {
	return *me.Address()
}

func (me *MpLockTm) Update(v uint64) {
	*me.Address() = v
}

func (me *MpLockTm) Init() {
	me.Update(0)
}

func (me *MpLockTm) LockR(timeout int) MpLockState {
	return me.until_ok(func() MpLockState {
		return me.lock(MPLOCK_READER, timeout)
	})
}

func (me *MpLockTm) LockW(timeout int) MpLockState {
	return me.until_ok(func() MpLockState {
		return me.lock(MPLOCK_WRITER, timeout)
	})
}

func (me *MpLockTm) UnLockR() MpLockState {
	return me.until_ok(func() MpLockState {
		return me.unlock(MPLOCK_READER)
	})
}

func (me *MpLockTm) UnLockW() MpLockState {
	return me.until_ok(func() MpLockState {
		return me.unlock(MPLOCK_WRITER)
	})
}

func (me *MpLockTm) Upgrade_rtow(timeout int) MpLockState {
	return me.until_ok(func() MpLockState {
		return me.upgrade(MPLOCK_READER, timeout)
	})
}

func (me *MpLockTm) Upgrade_wtor(timeout int) MpLockState {
	return me.until_ok(func() MpLockState {
		return me.upgrade(MPLOCK_WRITER, timeout)
	})
}

func (me *MpLockTm) Degrade_rtow(timeout int) MpLockState {
	return me.until_ok(func() MpLockState {
		return me.degrade(MPLOCK_READER, timeout)
	})
}

func (me *MpLockTm) Degrade_wtor(timeout int) MpLockState {
	return me.until_ok(func() MpLockState {
		return me.degrade(MPLOCK_WRITER, timeout)
	})
}

func (me *MpLockTm) HandleR(timeout int, handle func()) MpLockState {
	state := me.LockR(timeout)

	if MPLOCK_SECCESS == state {
		handle()

		me.UnLockR()
	}

	return state
}

func (me *MpLockTm) HandleW(timeout int, handle func()) MpLockState {
	state := me.LockW(timeout)

	if MPLOCK_SECCESS == state {
		handle()

		me.UnLockW()
	}

	return state
}
