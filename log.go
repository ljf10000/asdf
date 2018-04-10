package asdf

import (
	"fmt"
	"os"
	"time"
)

const (
	LogTypeInvalid LogType = 0
	LogTypeConsole LogType = 1
	LogTypeFile    LogType = 2
	LogTypeCo      LogType = 3
	LogTypeEnd     LogType = 4

	LogTypeDeft LogType = LogTypeConsole
)

type LogType int

var arrLogType = [LogTypeEnd]string{
	LogTypeInvalid: "invalid",
	LogTypeConsole: "console",
	LogTypeFile:    "file",
	LogTypeCo:      "co",
}
var mapLogType = map[string]LogType{
	"invalid": LogTypeInvalid,
	"console": LogTypeConsole,
	"file":    LogTypeFile,
	"co":      LogTypeCo,
}

func (me LogType) IsGood() bool {
	return me > LogTypeInvalid && me < LogTypeEnd
}

func (me LogType) String() string {
	if me.IsGood() {
		return arrLogType[me]
	} else {
		return Unknow
	}
}

func (me *LogType) FromString(s string) {
	if level, ok := mapLogType[s]; ok {
		*me = level
	}
}

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

func addLogPrefix(msg string) string {
	return time.Now().String() + ": " + msg
}

//==============================================================================

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

//==============================================================================

type FileLogger struct {
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

	logger := &FileLogger{
		level: LogLevelDeft,
		file:  file,
		fd:    fd,
		lock:  NewAccessLock("file-logger", false),
	}
	Log = logger

	return nil
}

func (me *FileLogger) GetLevel() LogLevel {
	return me.level
}

func (me *FileLogger) SetLevel(level LogLevel) {
	me.level = level
}

func (me *FileLogger) Log(level LogLevel, format string, v ...interface{}) {
	if level <= me.level {
		me.lock.Handle(func() {
			msg := addLogPrefix(fmt.Sprintf(format, v...))

			me.fd.WriteString(msg)
		})
	}
}

func (me *FileLogger) Emerg(format string, v ...interface{}) {
	me.Log(LogLevelEmerg, format+Crlf, v...)
}

func (me *FileLogger) Alert(format string, v ...interface{}) {
	me.Log(LogLevelAlert, format+Crlf, v...)
}

func (me *FileLogger) Crit(format string, v ...interface{}) {
	me.Log(LogLevelCrit, format+Crlf, v...)
}

func (me *FileLogger) Error(format string, v ...interface{}) {
	me.Log(LogLevelError, format+Crlf, v...)
}

func (me *FileLogger) Warning(format string, v ...interface{}) {
	me.Log(LogLevelWarning, format+Crlf, v...)
}

func (me *FileLogger) Notice(format string, v ...interface{}) {
	me.Log(LogLevelNotice, format+Crlf, v...)
}

func (me *FileLogger) Info(format string, v ...interface{}) {
	me.Log(LogLevelInfo, format+Crlf, v...)
}

func (me *FileLogger) Debug(format string, v ...interface{}) {
	me.Log(LogLevelDebug, format+Crlf, v...)
}

func (me *FileLogger) Close() error {
	return me.fd.Close()
}

//==============================================================================

type CoLogger struct {
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

	logger := &CoLogger{
		level: LogLevelDeft,
		file:  file,
		fd:    fd,
		ch:    ch,
	}
	Log = logger

	go logger.run(ch)

	return nil
}

func (me *CoLogger) run(ch chan string) {
	fmt.Printf("cologger running...\n")

	for {
		msg, ok := <-ch
		if !ok {
			me.fd.Close()

			return
		}

		msg = addLogPrefix(msg)
		me.fd.WriteString(msg)
		fmt.Print(msg)
	}
}

func (me *CoLogger) GetLevel() LogLevel {
	return me.level
}

func (me *CoLogger) SetLevel(level LogLevel) {
	me.level = level
}

func (me *CoLogger) Log(level LogLevel, format string, v ...interface{}) {
	if level <= me.level {
		me.ch <- fmt.Sprintf(format, v...)
	}
}

func (me *CoLogger) Emerg(format string, v ...interface{}) {
	me.Log(LogLevelEmerg, format+Crlf, v...)
}

func (me *CoLogger) Alert(format string, v ...interface{}) {
	me.Log(LogLevelAlert, format+Crlf, v...)
}

func (me *CoLogger) Crit(format string, v ...interface{}) {
	me.Log(LogLevelCrit, format+Crlf, v...)
}

func (me *CoLogger) Error(format string, v ...interface{}) {
	me.Log(LogLevelError, format+Crlf, v...)
}

func (me *CoLogger) Warning(format string, v ...interface{}) {
	me.Log(LogLevelWarning, format+Crlf, v...)
}

func (me *CoLogger) Notice(format string, v ...interface{}) {
	me.Log(LogLevelNotice, format+Crlf, v...)
}

func (me *CoLogger) Info(format string, v ...interface{}) {
	me.Log(LogLevelInfo, format+Crlf, v...)
}

func (me *CoLogger) Debug(format string, v ...interface{}) {
	me.Log(LogLevelDebug, format+Crlf, v...)
}

func (me *CoLogger) Close() error {
	ch := me.ch
	me.ch = make(chan string)
	close(ch)

	return nil
}

//==============================================================================
var Log ILogger = Console
