package asdf

import (
	"fmt"
	"os"
)

const (
	LogLevelEmerg   LogLevel = 0
	LogLevelAlert   LogLevel = 1
	LogLevelCrit    LogLevel = 2
	LogLevelError   LogLevel = 3
	LogLevelWarning LogLevel = 4
	LogLevelNotice  LogLevel = 5
	LogLevelInfo    LogLevel = 6
	LogLevelDebug   LogLevel = 7
	LogLevelEnd     LogLevel = 8
)

type LogLevel int

var LogLevels = [LogLevelEnd]string{
	LogLevelEmerg:   "Emerg",
	LogLevelAlert:   "Alert",
	LogLevelCrit:    "Crit",
	LogLevelError:   "Error",
	LogLevelWarning: "Warning",
	LogLevelNotice:  "Notice",
	LogLevelInfo:    "Info",
	LogLevelDebug:   "Debug",
}

func (me LogLevel) IsGood() bool {
	return me >= 0 && me < LogLevelEnd
}

func (me LogLevel) String() string {
	if me.IsGood() {
		return LogLevels[me]
	} else {
		return Unknow
	}
}

//==============================================================================

type consoleLogger struct {
	level LogLevel
}

var logConsole = &consoleLogger{
	level: LogLevelInfo,
}

func (me *consoleLogger) GetLevel() LogLevel {
	return me.level
}

func (me *consoleLogger) SetLevel(level LogLevel) {
	me.level = level
}

func (me *consoleLogger) Log(level LogLevel, format string, v ...interface{}) {
	if me.level <= level {
		fmt.Printf(format, v...)
	}
}

func (me *consoleLogger) Emerg(format string, v ...interface{}) {
	me.Log(LogLevelEmerg, format+Crlf, v...)
}

func (me *consoleLogger) Alert(format string, v ...interface{}) {
	me.Log(LogLevelAlert, format+Crlf, v...)
}

func (me *consoleLogger) Crit(format string, v ...interface{}) {
	me.Log(LogLevelCrit, format+Crlf, v...)
}

func (me *consoleLogger) Error(format string, v ...interface{}) {
	me.Log(LogLevelError, format+Crlf, v...)
}

func (me *consoleLogger) Warning(format string, v ...interface{}) {
	me.Log(LogLevelWarning, format+Crlf, v...)
}

func (me *consoleLogger) Notice(format string, v ...interface{}) {
	me.Log(LogLevelNotice, format+Crlf, v...)
}

func (me *consoleLogger) Info(format string, v ...interface{}) {
	me.Log(LogLevelInfo, format+Crlf, v...)
}

func (me *consoleLogger) Debug(format string, v ...interface{}) {
	me.Log(LogLevelDebug, format+Crlf, v...)
}

//==============================================================================

type fileLogger struct {
	level LogLevel

	filename string
	lock     *AccessLock
	file     *os.File
}

func newFileLogger(filename string) (*fileLogger, error) {
	file, err := os.OpenFile(filename, os.O_RDWR|os.O_CREATE, 0666)
	if nil != err {
		return nil, err
	}

	return &fileLogger{
		levelLogger: levelLogger{
			level: LogLevelInfo,
		},
		filename: filename,
		file:     file,
		lock:     NewAccessLock("file-logger", false),
	}, nil
}

func (me *fileLogger) GetLevel() LogLevel {
	return me.level
}

func (me *fileLogger) SetLevel(level LogLevel) {
	me.level = level
}

func (me *fileLogger) Log(level LogLevel, format string, v ...interface{}) {
	if me.level <= level {
		me.lock.Handle(func() {
			me.file.WriteString(fmt.Sprintf(format, v...))
		})
	}
}

func (me *fileLogger) Emerg(format string, v ...interface{}) {
	me.Log(LogLevelEmerg, format+Crlf, v...)
}

func (me *fileLogger) Alert(format string, v ...interface{}) {
	me.Log(LogLevelAlert, format+Crlf, v...)
}

func (me *fileLogger) Crit(format string, v ...interface{}) {
	me.Log(LogLevelCrit, format+Crlf, v...)
}

func (me *fileLogger) Error(format string, v ...interface{}) {
	me.Log(LogLevelError, format+Crlf, v...)
}

func (me *fileLogger) Warning(format string, v ...interface{}) {
	me.Log(LogLevelWarning, format+Crlf, v...)
}

func (me *fileLogger) Notice(format string, v ...interface{}) {
	me.Log(LogLevelNotice, format+Crlf, v...)
}

func (me *fileLogger) Info(format string, v ...interface{}) {
	me.Log(LogLevelInfo, format+Crlf, v...)
}

func (me *fileLogger) Debug(format string, v ...interface{}) {
	me.Log(LogLevelDebug, format+Crlf, v...)
}

//==============================================================================
var Log ILogger = logConsole

func SetLogger(r ILogger) {
	Log = r
}

func UseConsoleLogger() {
	Log = logConsole
}

func UseFileLogger(filename string) {
	if nil != Log {
		f, ok := Log.(*fileLogger)
		if ok {
			f.file.Close()
		}
	}

	Log, _ = newFileLogger(filename)
}
