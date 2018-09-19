package asdf

import (
	"bytes"
)

func CompareBytes(a, b []byte) int {
	return bytes.Compare(a, b)
}

func CompareString(a, b string) int {
	if a == b {
		return 0
	} else if a > b {
		return 1
	} else {
		return -1
	}
}

func CompareByte(a, b byte) int {
	if a == b {
		return 0
	} else if a > b {
		return 1
	} else {
		return -1
	}
}

func CompareUintptr(a, b uintptr) int {
	if a == b {
		return 0
	} else if a > b {
		return 1
	} else {
		return -1
	}
}

func CompareInt(a, b int) int {
	if a == b {
		return 0
	} else if a > b {
		return 1
	} else {
		return -1
	}
}

func CompareInt8(a, b int8) int {
	if a == b {
		return 0
	} else if a > b {
		return 1
	} else {
		return -1
	}
}

func CompareInt16(a, b int16) int {
	if a == b {
		return 0
	} else if a > b {
		return 1
	} else {
		return -1
	}
}

func CompareInt32(a, b int32) int {
	if a == b {
		return 0
	} else if a > b {
		return 1
	} else {
		return -1
	}
}

func CompareInt64(a, b int64) int {
	if a == b {
		return 0
	} else if a > b {
		return 1
	} else {
		return -1
	}
}

func CompareUint(a, b uint) int {
	if a == b {
		return 0
	} else if a > b {
		return 1
	} else {
		return -1
	}
}

func CompareUint8(a, b uint8) int {
	if a == b {
		return 0
	} else if a > b {
		return 1
	} else {
		return -1
	}
}

func CompareUint16(a, b uint16) int {
	if a == b {
		return 0
	} else if a > b {
		return 1
	} else {
		return -1
	}
}

func CompareUint32(a, b uint32) int {
	if a == b {
		return 0
	} else if a > b {
		return 1
	} else {
		return -1
	}
}

func CompareUint64(a, b uint64) int {
	if a == b {
		return 0
	} else if a > b {
		return 1
	} else {
		return -1
	}
}
