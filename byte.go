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

func CloneBytes(bin []byte) []byte {
	obj := make([]byte, len(bin))

	copy(obj, bin)

	return obj
}

func ByteReverse(bin []byte) {
	count := len(bin)

	if count > 1 {
		for i := 0; i < count/2; i++ {
			bin[i], bin[count-i] = bin[count-i], bin[i]
		}
	}
}

type Byte byte

func MakeByte(high, low int) Byte {
	return Byte((high << 4) | low)
}

func (me Byte) SetHigh(high int) Byte {
	low := me.Low()

	return MakeByte(high, low)
}

func (me Byte) SetLow(low int) Byte {
	high := me.High()

	return MakeByte(high, low)
}

func (me Byte) High() int {
	return int(me >> 4)
}

func (me Byte) Low() int {
	return int(me & 0xf)
}
