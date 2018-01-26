package asdf

type Pad4 struct {
	_ byte
	_ byte
	_ byte
	_ byte
}

type Pad8 struct {
	_ Pad4
	_ Pad4
}

type Pad16 struct {
	_ Pad8
	_ Pad8
}

type Pad32 struct {
	_ Pad16
	_ Pad16
}

type Pad64 struct {
	_ Pad32
	_ Pad32
}

type Pad128 struct {
	_ Pad64
	_ Pad64
}

type Pad256 struct {
	_ Pad128
	_ Pad128
}

type Pad512 struct {
	_ Pad256
	_ Pad256
}

type Pad1024 struct {
	_ Pad512
	_ Pad512
}

const (
	SizeofPad4    = 4
	SizeofPad8    = 8
	SizeofPad16   = 16
	SizeofPad32   = 32
	SizeofPad64   = 64
	SizeofPad128  = 128
	SizeofPad256  = 256
	SizeofPad512  = 512
	SizeofPad1024 = 1024
)
