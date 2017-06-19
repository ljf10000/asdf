package asdf

import (
	"errors"
	"os"
)

const (
	StdErrOk = iota

	StdErrError
	StdErrEmpty
	StdErrFull
	StdErrExist
	StdErrHolding
	StdErrPending
	StdErrTimeout

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
	StdErrBadPktLen
	StdErrBadPktDir
	StdErrBadVersion

	StdErrBadHttpClientPost

	StdErrNilObj
	StdErrNilBuffer
	StdErrNilIntf

	StdErrTooMore
	StdErrTooShortBuffer
	StdErrPktLenNoMatchBufferLen
)

var stdErrMap = map[int]string{
	StdErrOk: "ok",

	StdErrError:   "error",
	StdErrEmpty:   "empty",
	StdErrFull:    "full",
	StdErrExist:   "exist",
	StdErrHolding: "holding",
	StdErrPending: "pending",
	StdErrTimeout: "timeout",

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

	StdErrBadID:      "bad ID",
	StdErrBadObj:     "bad obj",
	StdErrBadBuf:     "bad buf",
	StdErrBadLen:     "bad length",
	StdErrBadIdx:     "bad idx",
	StdErrBadFsm:     "bad fsm",
	StdErrBadIntf:    "bad interface",
	StdErrBadType:    "bad type",
	StdErrBadMac:     "bad mac",
	StdErrBadName:    "bad name",
	StdErrBadJson:    "bad json",
	StdErrBadBody:    "bad body",
	StdErrBadConf:    "bad config",
	StdErrBadFile:    "bad File",
	StdErrBadPktLen:  "invalid packet length",
	StdErrBadPktDir:  "bad packet dir",
	StdErrBadVersion: "bad version",

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

	ErrError   = newerrors(StdErrError)
	ErrEmpty   = newerrors(StdErrEmpty)
	ErrFull    = newerrors(StdErrFull)
	ErrExist   = newerrors(StdErrExist)
	ErrHolding = newerrors(StdErrHolding)
	ErrPending = newerrors(StdErrPending)
	ErrTimeout = newerrors(StdErrTimeout)

	ErrNoPending   = newerrors(StdErrNoPending)
	ErrNoSupport   = newerrors(StdErrNoSupport)
	ErrNoExist     = newerrors(StdErrNoExist)
	ErrNoFound     = newerrors(StdErrNoFound)
	ErrNoFile      = newerrors(StdErrNoFile)
	ErrNoDir       = newerrors(StdErrNoDir)
	ErrNoMatch     = newerrors(StdErrNoMatch)
	ErrNoSpace     = newerrors(StdErrNoSpace)
	ErrNoPermit    = newerrors(StdErrNoPermit)
	ErrNoAlign     = newerrors(StdErrNoAlign)
	ErrNoIntf      = newerrors(StdErrNoIntf)
	ErrNoIpAddress = newerrors(StdErrNoIpAddress)
	ErrNoConnected = newerrors(StdErrNoConnected)

	ErrBadID      = newerrors(StdErrBadID)
	ErrBadObj     = newerrors(StdErrBadObj)
	ErrBadBuf     = newerrors(StdErrBadBuf)
	ErrBadLen     = newerrors(StdErrBadLen)
	ErrBadIdx     = newerrors(StdErrBadIdx)
	ErrBadFsm     = newerrors(StdErrBadFsm)
	ErrBadIntf    = newerrors(StdErrBadIntf)
	ErrBadType    = newerrors(StdErrBadType)
	ErrBadMac     = newerrors(StdErrBadMac)
	ErrBadName    = newerrors(StdErrBadName)
	ErrBadJson    = newerrors(StdErrBadJson)
	ErrBadBody    = newerrors(StdErrBadBody)
	ErrBadConf    = newerrors(StdErrBadConf)
	ErrBadFile    = newerrors(StdErrBadFile)
	ErrBadPktLen  = newerrors(StdErrBadPktLen)
	ErrBadPktDir  = newerrors(StdErrBadPktDir)
	ErrBadVersion = newerrors(StdErrBadVersion)

	ErrBadHttpClientPost = newerrors(StdErrBadHttpClientPost)

	ErrNilObj    = newerrors(StdErrNilObj)
	ErrNilBuffer = newerrors(StdErrNilBuffer)
	ErrNilIntf   = newerrors(StdErrNilIntf)

	ErrTooMore                = newerrors(StdErrTooMore)
	ErrTooShortBuffer         = newerrors(StdErrTooShortBuffer)
	ErrPktLenNoMatchBufferLen = newerrors(StdErrPktLenNoMatchBufferLen)
)

func newerrors(err int) error {
	return errors.New(StdErrorString(err))
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

func StdErrorString(error int) string {
	return stdErrMap[error]
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
