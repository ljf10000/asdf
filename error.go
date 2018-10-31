package asdf

import (
	"errors"
	"fmt"
	"os"
)

const (
	StdErrOk = iota

	StdErrError
	StdErrEmpty
	StdErrFull
	StdErrExist
	StdErrExistSame
	StdErrHolding
	StdErrPending
	StdErrTimeout
	StdErrListen
	StdErrStop
	StdErrLimit

	StdErrNoEnv
	StdErrNoPending
	StdErrNoSupport
	StdErrNoExist
	StdErrNoFound
	StdErrNoFile
	StdErrNoDir
	StdErrNoMatch
	StdErrNoSpace
	StdErrNoPermit
	StdErrNoAlign
	StdErrNoIntf
	StdErrNoIpAddress
	StdErrNoConnected

	StdErrBadID
	StdErrBadObj
	StdErrBadBuf
	StdErrBadLen
	StdErrBadIdx
	StdErrBadFsm
	StdErrBadIntf
	StdErrBadType
	StdErrBadMac
	StdErrBadName
	StdErrBadJson
	StdErrBadBody
	StdErrBadConf
	StdErrBadFile
	StdErrBadProto
	StdErrBadPktLen
	StdErrBadPktDir
	StdErrBadVersion
	StdErrBadIpAddress

	StdErrBadHttpClientPost

	StdErrNilObj
	StdErrNilBuffer
	StdErrNilIntf

	StdErrTooMore
	StdErrTooShortBuffer
	StdErrPktLenNoMatchBufferLen

	StdErrEnd
)

var stdErrMap = map[int]string{
	StdErrOk: "ok",

	StdErrError:     "error",
	StdErrEmpty:     "empty",
	StdErrFull:      "full",
	StdErrExist:     "exist",
	StdErrExistSame: "exist same",
	StdErrHolding:   "holding",
	StdErrPending:   "pending",
	StdErrTimeout:   "timeout",
	StdErrListen:    "listen",
	StdErrStop:      "stop",
	StdErrLimit:     "limit",

	StdErrNoEnv:       "no env",
	StdErrNoPending:   "no pending",
	StdErrNoSupport:   "no support",
	StdErrNoExist:     "no exist",
	StdErrNoFound:     "no found",
	StdErrNoFile:      "no file",
	StdErrNoDir:       "no dir",
	StdErrNoMatch:     "no match",
	StdErrNoSpace:     "no space",
	StdErrNoPermit:    "no permit",
	StdErrNoAlign:     "no align",
	StdErrNoIntf:      "no intf",
	StdErrNoIpAddress: "no ipaddress",
	StdErrNoConnected: "no connected",

	StdErrBadID:        "bad ID",
	StdErrBadObj:       "bad obj",
	StdErrBadBuf:       "bad buf",
	StdErrBadLen:       "bad length",
	StdErrBadIdx:       "bad idx",
	StdErrBadFsm:       "bad fsm",
	StdErrBadIntf:      "bad interface",
	StdErrBadType:      "bad type",
	StdErrBadMac:       "bad mac",
	StdErrBadName:      "bad name",
	StdErrBadJson:      "bad json",
	StdErrBadBody:      "bad body",
	StdErrBadConf:      "bad config",
	StdErrBadFile:      "bad file",
	StdErrBadProto:     "bad proto",
	StdErrBadPktLen:    "invalid packet length",
	StdErrBadPktDir:    "bad packet dir",
	StdErrBadVersion:   "bad version",
	StdErrBadIpAddress: "bad ip address",

	StdErrBadHttpClientPost: "http client post error",

	StdErrNilObj:    "nil obj",
	StdErrNilBuffer: "nil buffer",
	StdErrNilIntf:   "nil interface",

	StdErrTooMore:                "too more",
	StdErrTooShortBuffer:         "too short buffer",
	StdErrPktLenNoMatchBufferLen: "packet length not match buffer length",
}

var (
	Error = errors.New(Empty)

	ErrError     = NewError(StdErrError)
	ErrEmpty     = NewError(StdErrEmpty)
	ErrFull      = NewError(StdErrFull)
	ErrExist     = NewError(StdErrExist)
	ErrExistSame = NewError(StdErrExistSame)
	ErrHolding   = NewError(StdErrHolding)
	ErrPending   = NewError(StdErrPending)
	ErrTimeout   = NewError(StdErrTimeout)
	ErrListen    = NewError(StdErrListen)
	ErrStop      = NewError(StdErrStop)
	ErrLimit     = NewError(StdErrLimit)

	ErrNoEnv       = NewError(StdErrNoEnv)
	ErrNoPending   = NewError(StdErrNoPending)
	ErrNoSupport   = NewError(StdErrNoSupport)
	ErrNoExist     = NewError(StdErrNoExist)
	ErrNoFound     = NewError(StdErrNoFound)
	ErrNoFile      = NewError(StdErrNoFile)
	ErrNoDir       = NewError(StdErrNoDir)
	ErrNoMatch     = NewError(StdErrNoMatch)
	ErrNoSpace     = NewError(StdErrNoSpace)
	ErrNoPermit    = NewError(StdErrNoPermit)
	ErrNoAlign     = NewError(StdErrNoAlign)
	ErrNoIntf      = NewError(StdErrNoIntf)
	ErrNoIpAddress = NewError(StdErrNoIpAddress)
	ErrNoConnected = NewError(StdErrNoConnected)

	ErrBadID        = NewError(StdErrBadID)
	ErrBadObj       = NewError(StdErrBadObj)
	ErrBadBuf       = NewError(StdErrBadBuf)
	ErrBadLen       = NewError(StdErrBadLen)
	ErrBadIdx       = NewError(StdErrBadIdx)
	ErrBadFsm       = NewError(StdErrBadFsm)
	ErrBadIntf      = NewError(StdErrBadIntf)
	ErrBadType      = NewError(StdErrBadType)
	ErrBadMac       = NewError(StdErrBadMac)
	ErrBadName      = NewError(StdErrBadName)
	ErrBadJson      = NewError(StdErrBadJson)
	ErrBadBody      = NewError(StdErrBadBody)
	ErrBadConf      = NewError(StdErrBadConf)
	ErrBadFile      = NewError(StdErrBadFile)
	ErrBadProto     = NewError(StdErrBadProto)
	ErrBadPktLen    = NewError(StdErrBadPktLen)
	ErrBadPktDir    = NewError(StdErrBadPktDir)
	ErrBadVersion   = NewError(StdErrBadVersion)
	ErrBadIpAddress = NewError(StdErrBadIpAddress)

	ErrBadHttpClientPost = NewError(StdErrBadHttpClientPost)

	ErrNilObj    = NewError(StdErrNilObj)
	ErrNilBuffer = NewError(StdErrNilBuffer)
	ErrNilIntf   = NewError(StdErrNilIntf)

	ErrTooMore                = NewError(StdErrTooMore)
	ErrTooShortBuffer         = NewError(StdErrTooShortBuffer)
	ErrPktLenNoMatchBufferLen = NewError(StdErrPktLenNoMatchBufferLen)
)

func NewError(err int) error {
	return errors.New(StdErrorString(err))
}

func ErrSprintf(format string, v ...interface{}) error {
	return errors.New(fmt.Sprintf(format, v...))
}

func ErrLog(format string, v ...interface{}) error {
	s := fmt.Sprintf(format, v...)

	Log.Error("%s", s)

	return errors.New(s)
}

type StdError struct {
	Error       int    `json:"error"`
	ErrorString string `json:"errors"`
}

func (me *StdError) ObjOwner() string {
	return Unknow
}

func (me *StdError) ObjType() string {
	return "std-error"
}

func (me *StdError) ObjName() string {
	return me.ErrorString
}

func (me *StdError) ObjValue() string {
	return me.ErrorString
}

func StdErrorString(err int) string {
	if err >= 0 && err < StdErrEnd {
		return stdErrMap[err]
	} else {
		return "Unknow-Error"
	}
}

func NewStdError(error int, desc ...string) *StdError {
	if len(desc) > 0 {
		return &StdError{
			Error:       error,
			ErrorString: desc[0],
		}
	} else {
		return &StdError{
			Error:       error,
			ErrorString: StdErrorString(error),
		}
	}
}

func ExitError(error int) {
	Log.Error("exit error:%s", StdErrorString(error))

	os.Exit(error)
}

type MixError struct {
	StdError

	Err error
}

func (me *MixError) ObjOwner() string {
	return Unknow
}

func (me *MixError) ObjType() string {
	return "mix-error"
}

func (me *MixError) ObjName() string {
	if nil != me.Err {
		return me.Err.Error()
	} else {
		return me.ErrorString
	}
}

func (me *MixError) ObjValue() string {
	if nil != me.Err {
		return me.Err.Error()
	} else {
		return me.ErrorString
	}
}

func NewMixError(err error) *MixError {
	return &MixError{
		Err: err,
	}
}
