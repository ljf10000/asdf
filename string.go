package asdf

import (
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
