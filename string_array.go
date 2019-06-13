package asdf

import (
	"strings"
)

func NewStringArrayEx(n int) StringArray {
	return StringArray{
		ss: make([]string, n),
	}
}

func NewStringArray(n int) StringArray {
	const (
		min = 128
		max = 1024
	)

	if n < min {
		n = min
	} else if n > max {
		n = max
	}

	return NewStringArrayEx(n)
}

type StringArray struct {
	ss  []string
	cur int
}

func (me *StringArray) grow(n int) {
	count := len(me.ss)
	if count > n {
		n += count
	} else {
		n += n
	}

	ss := make([]string, n)

	copy(ss, me.ss)

	me.ss = ss
}

func (me *StringArray) Add(v ...string) {
	count := len(v)

	// BUGFIX: bug is <
	// 	fuck! fuck!! fuck!!!
	if me.cur+count > len(me.ss) {
		me.grow(count)
	}

	for i := 0; i < count; i++ {
		me.ss[me.cur+i] = v[i]

		// BUGFIX: out of memory
		// me.cur++
	}

	me.cur += count
}

func (me *StringArray) Build(sep string) string {
	if me.cur > 0 {
		return strings.Join(me.ss[:me.cur], sep)
	} else {
		return Empty
	}
}

func (me *StringArray) String() string {
	return me.Build(", ")
}
