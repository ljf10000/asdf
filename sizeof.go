package asdf

const (
	SizeofK = 1024
	SizeofM = 1024 * SizeofK
	SizeofG = 1024 * SizeofM
	SizeofT = 1024 * SizeofG
	SizeofP = 1024 * SizeofT
	SizeofE = 1024 * SizeofP

	SizeofByte    = 1
	SizeofInt8    = 1
	SizeofInt16   = 2
	SizeofInt32   = 4
	SizeofInt64   = 8
	SizeofFloat32 = 4
	SizeofFloat64 = 8
	SizeofPointer = 8

	SizeofPage      = 4 * SizeofK
	SizeofCacheLine = 64
)
