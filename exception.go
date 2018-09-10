package asdf

import (
	"fmt"
	"time"
)

func Panic(format string, v ...interface{}) {
	time.Sleep(3 * time.Second)

	panic(fmt.Sprintf(format, v...))
}
