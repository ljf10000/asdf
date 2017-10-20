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

func IsByteUpper(b byte) bool {
	return b >= 'A' && b <= 'Z'
}

func IsByteLower(b byte) bool {
	return b >= 'a' && b <= 'z'
}

func IsByteNumber(b byte) bool {
	return b >= '0' && b <= '9'
}

func IsByteUnder(b byte) bool {
	return b == '_'
}
