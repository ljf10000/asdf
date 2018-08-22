package asdf

import (
	"fmt"
)

type ConsoleLogger struct {
	level LogLevel
}

var Console = &ConsoleLogger{
	level: LogLevelDeft,
}

func (me *ConsoleLogger) GetLevel() LogLevel {
	return me.level
}

func (me *ConsoleLogger) SetLevel(level LogLevel) {
	me.level = level
}

func (me *ConsoleLogger) Log(level LogLevel, format string, v ...interface{}) {
	if level <= me.level {
		fmt.Printf(format, v...)
	}
}

func (me *ConsoleLogger) Emerg(format string, v ...interface{}) {
	me.Log(LogLevelEmerg, format+Crlf, v...)
}

func (me *ConsoleLogger) Alert(format string, v ...interface{}) {
	me.Log(LogLevelAlert, format+Crlf, v...)
}

func (me *ConsoleLogger) Crit(format string, v ...interface{}) {
	me.Log(LogLevelCrit, format+Crlf, v...)
}

func (me *ConsoleLogger) Error(format string, v ...interface{}) {
	me.Log(LogLevelError, format+Crlf, v...)
}

func (me *ConsoleLogger) Warning(format string, v ...interface{}) {
	me.Log(LogLevelWarning, format+Crlf, v...)
}

func (me *ConsoleLogger) Notice(format string, v ...interface{}) {
	me.Log(LogLevelNotice, format+Crlf, v...)
}

func (me *ConsoleLogger) Info(format string, v ...interface{}) {
	me.Log(LogLevelInfo, format+Crlf, v...)
}

func (me *ConsoleLogger) Debug(format string, v ...interface{}) {
	me.Log(LogLevelDebug, format+Crlf, v...)
}

func (me *ConsoleLogger) Close() error {
	return nil
}
