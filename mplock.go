package asdf

import (
	"time"
)

type MpLockState int

type mplock_role_t int
type mplock_can_t int
type mplock_act_t int

var (
	MpLockDump   = false
	MpLockUpTime = false
)

func mplock_time() Time32 {
	if MpLockUpTime {
		return NowUpTime()
	} else {
		return NowTime32()
	}
}

func MpLockPause() {
	time.Sleep(1)
}

const (
	MPLOCK_READER mplock_role_t = 0
	MPLOCK_WRITER mplock_role_t = 1
)

func (me mplock_role_t) String() string {
	if MPLOCK_READER == me {
		return "reader"
	} else {
		return "writer"
	}
}

const (
	MPLOCK_CAN_NONE mplock_can_t = 0
	MPLOCK_CAN_EXEC mplock_can_t = 1
	MPLOCK_CAN_WANT mplock_can_t = 2
)

func (me mplock_can_t) String() string {
	switch me {
	case MPLOCK_CAN_EXEC:
		return "exec"
	case MPLOCK_CAN_WANT:
		return "want"
	case MPLOCK_CAN_NONE:
		fallthrough
	default:
		return "none"
	}
}

const (
	MPLOCK_FAILED  MpLockState = 0
	MPLOCK_SECCESS MpLockState = 1
	MPLOCK_TIMEOUT MpLockState = 2
)

func (me MpLockState) String() string {
	switch me {
	case MPLOCK_TIMEOUT:
		return "timeout"
	case MPLOCK_SECCESS:
		return "seccess"
	case MPLOCK_FAILED:
		fallthrough
	default:
		return "failed"
	}
}

const (
	MPLOCK_ACT_LOCK    mplock_act_t = 0
	MPLOCK_ACT_UNLOCK  mplock_act_t = 1
	MPLOCK_ACT_UPGRADE mplock_act_t = 2
	MPLOCK_ACT_DEGRADE mplock_act_t = 3
	MPLOCK_ACT_END     mplock_act_t = 4
)

func (me mplock_act_t) String() string {
	switch me {
	case MPLOCK_ACT_LOCK:
		return "lock"
	case MPLOCK_ACT_UNLOCK:
		return "unlock"
	case MPLOCK_ACT_UPGRADE:
		return "upgrade"
	case MPLOCK_ACT_DEGRADE:
		fallthrough
	default:
		return "degrade"
	}
}

const (
	MPLOCK_VALUE_BITS = 15
	MPLOCK_VALUE_MASK = (1 << MPLOCK_VALUE_BITS) - 1
	MPLOCK_WANT_MASK  = 1 << MPLOCK_VALUE_BITS
)

type mplock struct {
	v uint16
}

func (me *mplock) init() {
	me.v = 0
}

func (me *mplock) want() bool {
	return MPLOCK_WANT_MASK == (me.v & MPLOCK_WANT_MASK)
}

func (me *mplock) value() uint16 {
	return me.v & MPLOCK_VALUE_MASK
}

func (me *mplock) locked() bool {
	return me.value() > 0
}

func (me *mplock) not_locked() bool {
	return me.value() == 0
}

func (me *mplock) multi_locked() bool {
	return me.value() > 1
}

func (me *mplock) set_want(want bool) {
	if want {
		me.v = me.value() | MPLOCK_WANT_MASK
	} else {
		me.v = me.value()
	}
}

func (me *mplock) set_value(value uint16) {
	me.v = (me.v & MPLOCK_WANT_MASK) | value
}

func (me *mplock) add() {
	me.v += 1
}

func (me *mplock) sub() {
	me.v -= 1
}

func mplock_check(role mplock_role_t, self *mplock, other *mplock) error {
	if self.want() && other.want() {
		return ErrLog("rw lock have double wanted")
	} else if self.locked() && other.locked() {
		return ErrLog("rw lock have double locked")
	} else if MPLOCK_READER == role && self.want() {
		return ErrLog("reader lock have wanted")
	} else if MPLOCK_WRITER == role && self.locked() {
		return ErrLog("writer lock have locked %d times", self.value())
	}

	return nil
}

type mplock_can_f func(role mplock_role_t, self *mplock, other *mplock) mplock_can_t

// return
//  MPLOCK_CAN_NONE:  can nothing
//  MPLOCK_CAN_EXEC:  can lock
//  MPLOCK_CAN_WANT:  can want
func mplock_can_lock(role mplock_role_t, self *mplock, other *mplock) mplock_can_t {
	//
	// self try lock
	//

	if other.locked() {
		// other have locked

		if MPLOCK_WRITER == role {
			//
			// self is writer, but other(reader) have locked the lock
			//
			if self.want() {
				return MPLOCK_CAN_NONE
			} else {
				return MPLOCK_CAN_WANT
			}
		} else {
			// self is reader
			return MPLOCK_CAN_NONE
		}
	} else if other.want() {
		// other NOT locked
		// other wanted the lock
		return MPLOCK_CAN_NONE
	} else if self.locked() {
		// other NOT locked
		// other NOT wanted
		//
		// self have locked
		//
		if MPLOCK_READER == role {
			// self is reader
			return MPLOCK_CAN_EXEC
		} else {
			// self is writer
			return MPLOCK_CAN_NONE
		}
	} else {
		// other NOT locked
		// other NOT wanted
		//
		// self is single lock, and NOT locked
		//  or
		// self is multi lock
		return MPLOCK_CAN_EXEC
	}
}

// return
//  MPLOCK_CAN_NONE:  can nothing
//  MPLOCK_CAN_EXEC:  can unlock
func mplock_can_unlock(role mplock_role_t, self *mplock, other *mplock) mplock_can_t {
	//
	// self try unlock(have locked)
	//

	if other.locked() {
		Panic("other have locked")
		return MPLOCK_CAN_NONE
	} else if self.value() == 0 {
		Panic("self  NOT locked")
		return MPLOCK_CAN_NONE
	} else {
		// other NOT  locked
		// self  have locked
		return MPLOCK_CAN_EXEC
	}
}

// return
//  MPLOCK_CAN_NONE:  can nothing
//  MPLOCK_CAN_EXEC:  can upgrade
func mplock_can_upgrade(role mplock_role_t, self *mplock, other *mplock) mplock_can_t {
	//
	// self try upgrade
	//
	if MPLOCK_WRITER == role {
		// self is writer
		// cann't upgrade to other
		Panic("self is writer")
		return MPLOCK_CAN_NONE
	} else if other.locked() {
		// other have locked
		Panic("other have locked")
		return MPLOCK_CAN_NONE
	} else if other.want() {
		// other have wanted
		return MPLOCK_CAN_NONE
	} else if self.not_locked() {
		// self NOT lock, cann't upgrade to other
		Panic("self NOT lock")
		return MPLOCK_CAN_NONE
	} else if self.multi_locked() {
		// self have locked multi
		return MPLOCK_CAN_NONE
	} else {
		return MPLOCK_CAN_EXEC
	}
}

// return
//  MPLOCK_CAN_NONE:  can nothing
//  MPLOCK_CAN_EXEC:  can degrade
func mplock_can_degrade(role mplock_role_t, self *mplock, other *mplock) mplock_can_t {
	//
	// self try degrade
	//

	if MPLOCK_READER == role {
		// self is reader
		// cann't degrade to ther
		Panic("self is reader")
		return MPLOCK_CAN_NONE
	} else if self.not_locked() {
		// self NOT locked
		Panic("self NOT locked")
		return MPLOCK_CAN_NONE
	} else if self.multi_locked() {
		// self have locked multi
		Panic("self have locked multi")
		return MPLOCK_CAN_NONE
	} else if other.locked() {
		// other have locked
		return MPLOCK_CAN_NONE
	} else {
		return MPLOCK_CAN_EXEC
	}
}
