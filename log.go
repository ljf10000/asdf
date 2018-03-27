package asdf

import (
	"fmt"
	"os"
)

const (
	LogLevelInvalid LogLevel = 0
	LogLevelEmerg   LogLevel = 1
	LogLevelAlert   LogLevel = 2
	LogLevelCrit    LogLevel = 3
	LogLevelError   LogLevel = 4
	LogLevelWarning LogLevel = 5
	LogLevelNotice  LogLevel = 6
	LogLevelInfo    LogLevel = 7
	LogLevelDebug   LogLevel = 8
	LogLevelEnd     LogLevel = 9

	LogLevelDeft LogLevel = LogLevelInfo
)

type LogLevel int

var arrLogLevel = [LogLevelEnd]string{
	LogLevelInvalid: "invalid",
	LogLevelEmerg:   "emerg",
	LogLevelAlert:   "alert",
	LogLevelCrit:    "crit",
	LogLevelError:   "error",
	LogLevelWarning: "warning",
	LogLevelNotice:  "notice",
	LogLevelInfo:    "info",
	LogLevelDebug:   "debug",
}

var mapLogLevel = map[string]LogLevel{
	"emerg":   LogLevelEmerg,
	"alert":   LogLevelAlert,
	"crit":    LogLevelCrit,
	"error":   LogLevelError,
	"warning": LogLevelWarning,
	"notice":  LogLevelNotice,
	"info":    LogLevelInfo,
	"debug":   LogLevelDebug,
}

func (me LogLevel) IsGood() bool {
	return me > LogLevelInvalid && me < LogLevelEnd
}

func (me LogLevel) String() string {
	if me.IsGood() {
		return arrLogLevel[me]
	} else {
		return Unknow
	}
}

func (me *LogLevel) FromString(s string) {
	if level, ok := mapLogLevel[s]; ok {
		*me = level
	}
}

//==============================================================================

type consoleLogger struct {
	level LogLevel
}

var logConsole = &consoleLogger{
	level: LogLevelDeft,
}

func (me *consoleLogger) GetLevel() LogLevel {
	return me.level
}

func (me *consoleLogger) SetLevel(level LogLevel) {
	me.level = level
}

func (me *consoleLogger) Log(level LogLevel, format string, v ...interface{}) {
	if level <= me.level {
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

func (me *consoleLogger) Close() error {
	return nil
}

//==============================================================================

type fileLogger struct {
	level LogLevel

	file string
	lock *AccessLock
	fd   *os.File
}

func OpenFileLogger(file string) error {
	fd, err := os.OpenFile(file, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if nil != err {
		return err
	}

	logger := &fileLogger{
		level: LogLevelDeft,
		file:  file,
		fd:    fd,
		lock:  NewAccessLock("file-logger", false),
	}
	Log = logger

	return nil
}

func (me *fileLogger) GetLevel() LogLevel {
	return me.level
}

func (me *fileLogger) SetLevel(level LogLevel) {
	me.level = level
}

func (me *fileLogger) Log(level LogLevel, format string, v ...interface{}) {
	if level <= me.level {
		me.lock.Handle(func() {
			me.fd.WriteString(fmt.Sprintf(format, v...))
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

func (me *fileLogger) Close() error {
	return me.fd.Close()
}

//==============================================================================

type coLogger struct {
	level LogLevel
	file  string
	fd    *os.File
	ch    chan string
}

func OpenCoLogger(file string, size int) error {
	fd, err := os.OpenFile(file, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if nil != err {
		return err
	}

	ch := make(chan string, size)

	logger := &coLogger{
		level: LogLevelDeft,
		file:  file,
		fd:    fd,
		ch:    ch,
	}
	Log = logger

	go logger.run(ch)

	return nil
}

func (me *coLogger) run(ch chan string) {
	fmt.Printf("cologger running...\n")

	for {
		msg, ok := <-ch
		if !ok {
			me.fd.Close()

			return
		}

		me.fd.WriteString(msg)
		fmt.Print(msg)
	}
}

func (me *coLogger) GetLevel() LogLevel {
	return me.level
}

func (me *coLogger) SetLevel(level LogLevel) {
	me.level = level
}

func (me *coLogger) Log(level LogLevel, format string, v ...interface{}) {
	if level <= me.level {
		me.ch <- fmt.Sprintf(format, v...)
	}
}

func (me *coLogger) Emerg(format string, v ...interface{}) {
	me.Log(LogLevelEmerg, format+Crlf, v...)
}

func (me *coLogger) Alert(format string, v ...interface{}) {
	me.Log(LogLevelAlert, format+Crlf, v...)
}

func (me *coLogger) Crit(format string, v ...interface{}) {
	me.Log(LogLevelCrit, format+Crlf, v...)
}

func (me *coLogger) Error(format string, v ...interface{}) {
	me.Log(LogLevelError, format+Crlf, v...)
}

func (me *coLogger) Warning(format string, v ...interface{}) {
	me.Log(LogLevelWarning, format+Crlf, v...)
}

func (me *coLogger) Notice(format string, v ...interface{}) {
	me.Log(LogLevelNotice, format+Crlf, v...)
}

func (me *coLogger) Info(format string, v ...interface{}) {
	me.Log(LogLevelInfo, format+Crlf, v...)
}

func (me *coLogger) Debug(format string, v ...interface{}) {
	me.Log(LogLevelDebug, format+Crlf, v...)
}

func (me *coLogger) Close() error {
	ch := me.ch
	me.ch = make(chan string)
	close(ch)

	return nil
}

//==============================================================================
var Log ILogger = logConsole
