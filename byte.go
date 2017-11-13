package asdf

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

func GetBytes(bin []byte, offset, size int) ([]byte, int) {
	end := offset + size

	return bin[offset:end], end
}

func ByteReverse(bin []byte) {
	count := len(bin)

	if count > 1 {
		for i := 0; i < count/2; i++ {
			bin[i], bin[count-i] = bin[count-i], bin[i]
		}
	}
}
