package asdf

func Panicf(format string, v ...interface{}) {
	Log.Panic(format, v...)
}

func Panic(format string, v ...interface{}) {
	Log.Panic(format, v...)
}
