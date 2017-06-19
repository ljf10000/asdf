package asdf

import (
	"fmt"
)

////////////////////////////////////////////////////////////////////////////////
// single interface
////////////////////////////////////////////////////////////////////////////////
type IBegin interface {
	Begin() int
}

type IEnd interface {
	End() int
}

type IInt interface {
	Int() int
}

type IFloat interface {
	Float() float64
}

type ISlice interface {
	Slice() []byte
}

type IEq interface {
	Eq(interface{}) bool
}

type IGt interface {
	Gt(interface{}) bool
}

type ILt interface {
	Lt(interface{}) bool
}

type IFirst interface {
	First() interface{}
}

type ILast interface {
	Last() interface{}
}

type ITails interface {
	Tail() []interface{}
}

type IHead interface {
	Head() []interface{}
}

type IPrev interface {
	Prev() interface{}
}

type INext interface {
	Next() interface{}
}

type IReverse interface {
	Reverse() []interface{}
}

type IRepeat interface {
	Repeat(int) []interface{}
}

type IString interface {
	fmt.Stringer
}

type IFromString interface {
	FromString(string) error
}

type IGood interface {
	IsGood() bool
}

type INew interface {
	New() interface{}
}

type IFind interface {
	Find() (interface{}, bool)
}

type ISave interface {
	Save()
}

type IUnSave interface {
	UnSave()
}

////////////////////////////////////////////////////////////////////////////////
// combination interface
////////////////////////////////////////////////////////////////////////////////
type IBound interface {
	// [begin, end)
	IBegin
	IEnd
}

type INumber interface {
	IBound
	IInt
}

type ICompare interface {
	IEq
	IGt
}

type IList interface {
	IFirst
	ILast

	ITails
	IHead

	IReverse
}

type IListNode interface {
	IPrev
	INext
}

type IStorage interface {
	INew
	IFind
	ISave
	IUnSave
}

type ILogEmerg interface {
	Emerg(format string, v ...interface{})
}

type ILogAlert interface {
	Alert(format string, v ...interface{})
}

type ILogCrit interface {
	Crit(format string, v ...interface{})
}

type ILogError interface {
	Error(format string, v ...interface{})
}

type ILogWarning interface {
	Warning(format string, v ...interface{})
}

type ILogNotice interface {
	Notice(format string, v ...interface{})
}

type ILogInfo interface {
	Info(format string, v ...interface{})
}

type ILogDebug interface {
	Debug(format string, v ...interface{})
}

type ILogger interface {
	ILogEmerg
	ILogAlert
	ILogCrit
	ILogError
	ILogWarning
	ILogNotice
	ILogInfo
	ILogDebug
}

type IObjOwner interface {
	ObjOwner() string
}

type IObjType interface {
	ObjType() string
}

type IObjName interface {
	ObjName() string
}

type IObjValue interface {
	ObjValue() string
}

type IObj interface {
	IObjOwner
	IObjType
	IObjName
	//	IObjValue
}

type IEncode interface {
	Encode([]byte) []byte
}

type IDecode interface {
	Decode([]byte) ([]byte, error)
}

type ICodec interface {
	IEncode
	IDecode
}

type ISerialize interface {
	Serialize() error
}

type IUnserialize interface {
	Unserialize() error
}

type IFileName interface {
	FileName() FileName
}
