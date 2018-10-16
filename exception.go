package asdf

import (
	"fmt"
	"time"
)

func Panic(format string, v ...interface{}) {
	Log.Emerg("Panic: "+format, v...)
	time.Sleep(3 * time.Second)
	Log.Close()

	panic(fmt.Sprintf(format, v...))
}
