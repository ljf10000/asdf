package asdf

import (
	"fmt"
	"time"
)

func Panic(format string, v ...interface{}) {
	Log.Emerg("Panic: "+format, v...)

	for i := 0; i < 2000; i++ {
		time.Sleep(time.Millisecond)
	}

	Log.Close()

	panic(fmt.Sprintf(format, v...))
}
