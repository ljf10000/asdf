package asdf

import (
	"fmt"
)

/*
*   raw format(like UltraEdit)

      :                                     ;
 Line :       Hexadecimal Content           ; Raw Content
      : 0 1 2 3  4 5 6 7  8 9 A B  C D E F  ;
      :                                     ;
xxxxH : xxxxxxxx xxxxxxxx xxxxxxxx xxxxxxxx ; cccccccccccccccc
xxxxH : xxxxxxxx xxxxxxxx xxxxxxxx xxxxxxxx ; cccccccccccccccc
xxxxH : xxxxxxxx xxxxxxxx xxxxxxxx xxxxxxxx ; cccccccccccccccc
xxxxH : xxxxxxxx xxxxxxxx xxxxxxxx xxxxxxxx ; cccccccccccccccc
xxxxH : xxxxxxxx xxxxxxxx xxxxxxxx xxxxxxxx ; cccccccccccccccc
xxxxH : xxxxxxxx xxxxxxxx xxxxxx            ; ccccccccccc
*/

const (
	__DUMP_LINE_BLOCK       = 4
	__DUMP_LINE_BLOCK_BYTES = 4
	__DUMP_LINE_LIMIT       = 80

	__DUMP_LINE_BYTES = __DUMP_LINE_BLOCK_BYTES * __DUMP_LINE_BLOCK
	__DUMP_LINE_MAX   = (0 +
		8 + // "xxxxH : "
		(2*__DUMP_LINE_BLOCK_BYTES+1)*__DUMP_LINE_BLOCK +
		2 + // "; "
		__DUMP_LINE_BYTES + // BytesofDumpLine
		1) // "\n"

	__DUMP_LINE_SEPARATOR     = "=============================================================="
	__DUMP_LINE_SEPARATOR_SUB = "--------------------------------------------------------------"

	__DUMP_LINE_HEADER__ = `
      :                                     ;
 Line :       Hexadecimal Content           ; Raw Content
      : 0 1 2 3  4 5 6 7  8 9 A B  C D E F  ;
      :                                     ;
`
)

var __DUMP_LINE_HEADER = __DUMP_LINE_HEADER__[1:]

func lineSprintf(iLine int, bin []byte) string {
	Len := len(bin)

	/*
	 * line as
	 *
	 * "xxxxH :"
	 */
	s := fmt.Sprintf("%.4XH :", __DUMP_LINE_BYTES*iLine)

	/*
	 * Hexadecimal Content as
	 *
	 * " xxxxxxxx xxxxxxxx xxxxxxxx xxxxxxxx ; "
	 */
	for i := 0; i < Len; i++ {
		if 0 == (i % __DUMP_LINE_BLOCK_BYTES) {
			s += Space
		}

		s += fmt.Sprintf("%.2X", bin[i])
	}

	for i := Len; i < __DUMP_LINE_BYTES; i++ {
		if 0 == (i % __DUMP_LINE_BLOCK_BYTES) {
			s += Space
		}

		s += Space + Space
	}
	s += " ; "

	/*
	 * Raw Content as
	 *
	 * "cccccccccccccccc"
	 */
	for i := 0; i < Len; i++ {
		v := bin[i]

		if v >= 0x20 && v <= 0x7e {
			s += string(v)
		} else {
			s += FileNameSplit
		}
	}
	s += Crlf

	return s
}

func BinSprintf(bin []byte) string {
	Len := len(bin)
	if 0 == Len {
		return Empty
	}

	lineCount := AlignI(Len, __DUMP_LINE_BYTES) / __DUMP_LINE_BYTES

	tail := Len % __DUMP_LINE_BYTES
	if 0 == tail {
		tail = __DUMP_LINE_BYTES
	}

	// fmt.Printf("dump bin:%d\n", len(bin))

	s := __DUMP_LINE_HEADER
	for i := 0; i < lineCount-1; i++ {
		begin := i * __DUMP_LINE_BYTES
		end := begin + __DUMP_LINE_BYTES

		// fmt.Printf("dump line[%d] begin:%d end:%d\n", i, begin, end)
		s += lineSprintf(i, bin[begin:end])
	}

	begin := (lineCount - 1) * __DUMP_LINE_BYTES
	end := begin + tail

	// fmt.Printf("dump line[%d] begin:%d end:%d\n", lineCount-1, begin, end)
	s += lineSprintf(lineCount-1, bin[begin:end])

	return s
}
