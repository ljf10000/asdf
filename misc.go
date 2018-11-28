package asdf

const (
	READ_ONLY  RwMode = true
	READ_WRITE RwMode = false
)

type RwMode bool

func (me RwMode) String() string {
	if me {
		return "ro"
	} else {
		return "rw"
	}
}
