package asdf

const (
	Empty = ""
	Space = " "
	Tab   = Space + Space + Space + Space
	Tab2  = Tab + Tab
	Tab3  = Tab2 + Tab
	Tab4  = Tab3 + Tab
	Crlf  = "\n"
	Crlf2 = Crlf + Crlf
	Crlf3 = Crlf2 + Crlf
	Crlf4 = Crlf3 + Crlf

	Unknow = "unknow"
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
