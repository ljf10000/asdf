package asdf

import (
	"reflect"
	"strconv"
	"strings"
	"unicode/utf8"
	"unsafe"
)

const (
	Empty  = ""
	Space  = " "
	Unknow = "unknow"

	Space2 = Space + Space
	Space4 = Space2 + Space2
	Space8 = Space4 + Space4

	Tab  = Space + Space + Space + Space
	Tab2 = Tab + Tab
	Tab3 = Tab2 + Tab
	Tab4 = Tab3 + Tab

	TabR  = "\t"
	TabR2 = TabR + TabR
	TabR3 = TabR2 + TabR
	TabR4 = TabR3 + TabR

	Crlf  = "\n"
	Crlf2 = Crlf + Crlf
	Crlf3 = Crlf2 + Crlf
	Crlf4 = Crlf3 + Crlf
)

func SkipString(s string, skip int) string {
	if len(s) > skip {
		return s[skip:]
	} else {
		return Empty
	}
}

func SquareString(s string) string {
	return "[" + s + "]"
}

func CurvesString(s string) string {
	return "(" + s + ")"
}

func YesNo(yes bool) string {
	if yes {
		return "yes"
	} else {
		return "no"
	}
}

func OnOff(on bool) string {
	if on {
		return "on"
	} else {
		return "off"
	}
}

func RepeatN(r string, n int) string {
	s := Empty

	for i := 0; i < n; i++ {
		s += r
	}

	return s
}

func SpaceN(n int) string {
	return RepeatN(Space, n)
}

func TabN(n int) string {
	return RepeatN(Tab, n)
}

func CrlfN(n int) string {
	return RepeatN(Crlf, n)
}

func HasPrefix(s string, ss []string) (int, bool) {
	for k, v := range ss {
		if strings.HasPrefix(s, v) {
			return k, true
		}
	}

	return 0, false
}

func FirstRune(line string) (rune, string) {
	if len(line) > 0 {
		c, _ := utf8.DecodeRuneInString(line)

		return c, line[len(string(c)):]
	} else {
		return 0, Empty
	}
}

// Join concatenates the elements of a to create a single string. The separator string
// sep is placed between elements in the resulting string.
func StringsToBin(a []string, sep string) []byte {
	switch len(a) {
	case 0:
		return []byte("")
	case 1:
		return []byte(a[0])
	case 2:
		// Special case for common small values.
		// Remove if golang.org/issue/6714 is fixed
		return []byte(a[0] + sep + a[1])
	case 3:
		// Special case for common small values.
		// Remove if golang.org/issue/6714 is fixed
		return []byte(a[0] + sep + a[1] + sep + a[2])
	case 4:
		// Special case for common small values.
		// Remove if golang.org/issue/6714 is fixed
		return []byte(a[0] + sep + a[1] + sep + a[2] + sep + a[3])
	}

	n := len(sep) * (len(a) - 1)
	for i := 0; i < len(a); i++ {
		n += len(a[i])
	}

	b := make([]byte, n, n+1) // 1 for append crlf
	bp := copy(b, a[0])
	for _, s := range a[1:] {
		bp += copy(b[bp:], sep)
		bp += copy(b[bp:], s)
	}

	return b
}

type IField interface {
	Name() string
	String() string
}

func MakeFieldListString(fields ...IField) string {
	s := Empty

	for _, field := range fields {
		s += ", " + field.Name() + ":" + field.String()
	}

	return SkipString(s, 2)
}

/******************************************************************************/

func Utox(v uint) string {
	return strconv.FormatUint(uint64(v), 16)
}

func Utox8(v uint8) string {
	return strconv.FormatUint(uint64(v), 16)
}

func Utox16(v uint16) string {
	return strconv.FormatUint(uint64(v), 16)
}

func Utox32(v uint32) string {
	return strconv.FormatUint(uint64(v), 16)
}

func Utox64(v uint64) string {
	return strconv.FormatUint(uint64(v), 16)
}

/******************************************************************************/

func Utoa(v uint) string {
	return strconv.FormatUint(uint64(v), 10)
}

func Utoa8(v uint8) string {
	return strconv.FormatUint(uint64(v), 10)
}

func Utoa16(v uint16) string {
	return strconv.FormatUint(uint64(v), 10)
}

func Utoa32(v uint32) string {
	return strconv.FormatUint(uint64(v), 10)
}

func Utoa64(v uint64) string {
	return strconv.FormatUint(uint64(v), 10)
}

/******************************************************************************/

func Itox(v int) string {
	return strconv.FormatInt(int64(v), 16)
}

func Itox8(v int8) string {
	return strconv.FormatInt(int64(v), 16)
}

func Itox16(v int16) string {
	return strconv.FormatInt(int64(v), 16)
}

func Itox32(v int32) string {
	return strconv.FormatInt(int64(v), 16)
}

func Itox64(v int64) string {
	return strconv.FormatInt(int64(v), 16)
}

/******************************************************************************/

func Itoa(v int) string {
	return strconv.FormatInt(int64(v), 10)
}

func Itoa8(v int8) string {
	return strconv.FormatInt(int64(v), 10)
}

func Itoa16(v int16) string {
	return strconv.FormatInt(int64(v), 10)
}

func Itoa32(v int32) string {
	return strconv.FormatInt(int64(v), 10)
}

func Itoa64(v int64) string {
	return strconv.FormatInt(int64(v), 10)
}

/******************************************************************************/

func Atoi(s string) int {
	v, err := strconv.Atoi(s)
	if nil != err {
		return 0
	} else {
		return v
	}
}

/******************************************************************************/

type String string

func (me *String) Header() *reflect.StringHeader {
	return (*reflect.StringHeader)(unsafe.Pointer(me))
}

func (me String) Slice() []byte {
	h := me.Header()

	return MakeSliceEx(h.Data, h.Len, h.Len)
}

func (me String) Bin() []byte {
	return me.Slice()
}

func (me String) String() string {
	return string(me)
}

func (me String) Eq(it interface{}) bool {
	v, ok := it.(IString)
	if !ok {
		return false
	}

	return string(me) == v.String()
}

func ObjToString(obj unsafe.Pointer, size int) string {
	if nil != obj {
		var s String

		hdr := s.Header()
		hdr.Data = uintptr(obj)
		hdr.Len = size

		return string(s)
	} else {
		return Empty
	}
}

func MemberToString(obj unsafe.Pointer, offset, size uintptr) string {
	member := unsafe.Pointer(uintptr(obj) + offset)

	return ObjToString(member, int(size))
}

func MakeStringEx(Data uintptr, Len int) string {
	var s String

	h := s.Header()

	h.Data = Data
	h.Len = Len

	return string(s)
}

func MakeString(obj unsafe.Pointer, Len int) string {
	return MakeStringEx(uintptr(obj), Len)
}

func StringAddress(s string) uintptr {
	return ((*reflect.StringHeader)(unsafe.Pointer(&s))).Data
}

func StringPointer(s string) unsafe.Pointer {
	return unsafe.Pointer(((*reflect.StringHeader)(unsafe.Pointer(&s))).Data)
}

func StringToBin(v string) []byte {
	s := String(v)

	return s.Bin()
}

func BinToString(v []byte) string {
	return MakeString(SlicePointer(v), len(v))
}
