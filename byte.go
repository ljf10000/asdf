package asdf

func Copy(dst, src []byte) int {
	if nil != dst {
		if nil != src {
			return copy(dst, src)
		} else {
			return 0
		}
	} else {
		return 0
	}
}

func Len(buf []byte) int {
	if nil != buf {
		return len(buf)
	} else {
		return 0
	}
}
