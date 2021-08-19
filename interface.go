package asdf

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

type IEq interface {
	Eq(interface{}) bool
}

type IGt interface {
	Gt(interface{}) bool
}

type ILt interface {
	Lt(interface{}) bool
}

type IRepeat interface {
	Repeat(int) []interface{}
}

type IString interface {
	String() string
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

type IGet interface {
	Get() interface{}
}

type ISet interface {
	Set(v interface{})
}

type ISize interface {
	Size() int
}

type ICount interface {
	Count() int
}

type ISlice interface {
	Slice() []byte
}

type IToBinary interface {
	ToBinary(bin []byte) error
}

type IFromBinary interface {
	FromBinary(bin []byte) error
}

type IClose interface {
	Close() error
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

type IEncode interface {
	Encode([]byte) []byte
}

type IDecode interface {
	Decode([]byte) ([]byte, error)
}

type ISerialize interface {
	Serialize() error
}

type IUnserialize interface {
	Unserialize() error
}

type IMarshal interface {
	Marshal() ([]byte, error)
}

type IUnmarshal interface {
	Unmarshal(data []byte) error
}

type IFileName interface {
	FileName() FileName
}

type ISort interface {
	Len() int
	Less(i, j int) bool
	Swap(i, j int)
}

type ITimeout interface {
	Error() string
	Timeout() bool
	Temporary() bool
}

func IsTimeout(e interface{}) bool {
	_, ok := e.(ITimeout)

	return ok
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

type IStorage interface {
	INew
	IFind
	ISave
	IUnSave
}

type IObj interface {
	IObjOwner
	IObjType
	IObjName
	//	IObjValue
}

type ICodec interface {
	IEncode
	IDecode
}

type IBinary interface {
	ISize
	IToBinary
	IFromBinary
}

func ToBinary(obj IBinary) ([]byte, error) {
	bin := make([]byte, obj.Size())

	err := obj.ToBinary(bin)
	if nil != err {
		return nil, err
	} else {
		return bin, nil
	}
}

const (
	OpCompareEq  OpCompare = 0
	OpCompareGe  OpCompare = 1
	OpCompareLe  OpCompare = 2
	OpCompareEnd OpCompare = 3
)

type OpCompare int
