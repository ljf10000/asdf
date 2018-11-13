package asdf

import (
	"strconv"
	"strings"
	"unicode/utf8"
)

const (
	Empty  = ""
	Space  = " "
	Unknow = "unknow"

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

type IField interface {
	Name() string
	String() string
}

func PrefixString(s string) string {
	if len(s) > 0 {
		return s[2:]
	} else {
		return Empty
	}
}

func MakeFieldListString(fields ...IField) string {
	s := Empty

	for _, field := range fields {
		s += ", " + field.Name() + ":" + field.String()
	}

	return PrefixString(s)
}

func Utoa(v uint) string {
	return strconv.FormatUint(uint64(v), 10)
}

func UtoaX(v uint) string {
	return strconv.FormatUint(uint64(v), 16)
}

func UtoaO(v uint) string {
	return strconv.FormatUint(uint64(v), 8)
}

func Itoa(v int) string {
	return strconv.FormatInt(int64(v), 10)
}

func ItoaX(v int) string {
	return strconv.FormatInt(int64(v), 16)
}

func ItoaO(v int) string {
	return strconv.FormatInt(int64(v), 8)
}

func Atoi(s string) int {
	v, err := strconv.Atoi(s)
	if nil != err {
		return 0
	} else {
		return v
	}
}
