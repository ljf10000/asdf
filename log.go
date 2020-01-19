package asdf

import (
	"fmt"
	"time"
)

type ILogGetLevel interface {
	GetLevel() LogLevel
}

type ILogSetLevel interface {
	SetLevel(level LogLevel)
}

type ILogPanic interface {
	Panic(format string, v ...interface{})
}

type ILogLog interface {
	Log(level LogLevel, format string, v ...interface{})
}

type ILogEmerg interface {
	Emerg(format string, v ...interface{})
}

type ILogAlert interface {
	Alert(format string, v ...interface{})
}

type ILogCrit interface {
	Crit(format string, v ...interface{})
}

type ILogError interface {
	Error(format string, v ...interface{})
}

type ILogWarning interface {
	Warning(format string, v ...interface{})
}

type ILogWarn interface {
	Warn(format string, v ...interface{})
}

type ILogNotice interface {
	Notice(format string, v ...interface{})
}

type ILogInfo interface {
	Info(format string, v ...interface{})
}

type ILogDebug interface {
	Debug(format string, v ...interface{})
}

type ILogger interface {
	ILogEmerg
	ILogAlert
	ILogCrit
	ILogError
	ILogWarning
	ILogWarn
	ILogNotice
	ILogInfo
	ILogDebug

	ILogLog
	ILogPanic
	ILogGetLevel
	ILogSetLevel

	IClose
}

const (
	LogTypeInvalid LogType = 0
	LogTypeConsole LogType = 1
	LogTypeFile    LogType = 2
	LogTypeCo      LogType = 3
	LogTypeEnd     LogType = 4

	LogTypeDeft LogType = LogTypeConsole
)

type LogType int

var logTypes = EnumMapper{
	Type: "asdf.LogType",
	Names: []string{
		LogTypeInvalid: "invalid",
		LogTypeConsole: "console",
		LogTypeFile:    "file",
		LogTypeCo:      "co",
	},
}

func (me LogType) IsGood() bool {
	return me > LogTypeInvalid && me < LogTypeEnd
}

func (me LogType) String() string {
	if me.IsGood() {
		return logTypes.Name(int(me))
	} else {
		return Unknow
	}
}

func (me *LogType) FromString(s string) error {
	idx, err := logTypes.Index(s)
	if nil == err {
		*me = LogType(idx)
	}

	return err
}

const (
	LogLevelInvalid LogLevel = 0
	LogLevelEmerg   LogLevel = 1
	LogLevelAlert   LogLevel = 2
	LogLevelCrit    LogLevel = 3
	LogLevelError   LogLevel = 4
	LogLevelWarning LogLevel = 5
	LogLevelWarn    LogLevel = LogLevelWarning
	LogLevelNotice  LogLevel = 6
	LogLevelInfo    LogLevel = 7
	LogLevelDebug   LogLevel = 8
	LogLevelEnd     LogLevel = 9

	LogLevelDeft LogLevel = LogLevelInfo
)

type LogLevel byte

type LogLevelMapper struct {
	Name  string
	Short string
}

var logLevels = [LogLevelEnd]LogLevelMapper{
	LogLevelEmerg: LogLevelMapper{
		Name:  "emerg",
		Short: "M",
	},
	LogLevelAlert: LogLevelMapper{
		Name:  "alert",
		Short: "A",
	},
	LogLevelCrit: LogLevelMapper{
		Name:  "crit",
		Short: "C",
	},
	LogLevelError: LogLevelMapper{
		Name:  "error",
		Short: "E",
	},
	LogLevelWarning: LogLevelMapper{
		Name:  "warn",
		Short: "W",
	},
	LogLevelNotice: LogLevelMapper{
		Name:  "notice",
		Short: "N",
	},
	LogLevelInfo: LogLevelMapper{
		Name:  "info",
		Short: "I",
	},
	LogLevelDebug: LogLevelMapper{
		Name:  "debug",
		Short: "D",
	},
}

var mapLogLevel = map[string]LogLevel{
	"emerg":   LogLevelEmerg,
	"alert":   LogLevelAlert,
	"crit":    LogLevelCrit,
	"error":   LogLevelError,
	"warning": LogLevelWarning,
	"warn":    LogLevelWarning,
	"notice":  LogLevelNotice,
	"info":    LogLevelInfo,
	"debug":   LogLevelDebug,
}

func (me LogLevel) IsGood() bool {
	return me > LogLevelInvalid && me < LogLevelEnd
}

func (me LogLevel) String() string {
	if me.IsGood() {
		return logLevels[me].Name
	} else {
		return Unknow
	}
}

func (me LogLevel) Short() string {
	if me.IsGood() {
		return logLevels[me].Short
	} else {
		return "U"
	}
}

func (me *LogLevel) FromString(s string) {
	if level, ok := mapLogLevel[s]; ok {
		*me = level
	} else {
		*me = LogLevelInfo
	}
}

const LogTimeFormat = "2006-01-02@15:04:05.999999999"

type LogMsg struct {
	s     string
	level LogLevel
	panic bool
}

func (me *LogMsg) String() string {
	return fmt.Sprintf("[%s] %s %s",
		me.level.Short(),
		time.Now().Format(LogTimeFormat),
		me.s)
}

var Log ILogger = Console

/******************************************************************************/

type LogModule uint64

func (me *LogModule) isLog(module int) bool {
	flag := uint64(1) << uint64(module)

	return nil == me || flag == (flag&uint64(*me))
}

func (me *LogModule) Emerg(module int, format string, v ...interface{}) {
	if me.isLog(module) {
		Log.Log(LogLevelEmerg, format, v...)
	}
}

func (me *LogModule) Alert(module int, format string, v ...interface{}) {
	if me.isLog(module) {
		Log.Log(LogLevelAlert, format, v...)
	}
}

func (me *LogModule) Crit(module int, format string, v ...interface{}) {
	if me.isLog(module) {
		Log.Log(LogLevelCrit, format, v...)
	}
}

func (me *LogModule) Error(module int, format string, v ...interface{}) {
	if me.isLog(module) {
		Log.Log(LogLevelError, format, v...)
	}
}

func (me *LogModule) Warning(module int, format string, v ...interface{}) {
	if me.isLog(module) {
		Log.Log(LogLevelWarning, format, v...)
	}
}

func (me *LogModule) Warn(module int, format string, v ...interface{}) {
	if me.isLog(module) {
		Log.Log(LogLevelWarn, format, v...)
	}
}

func (me *LogModule) Notice(module int, format string, v ...interface{}) {
	if me.isLog(module) {
		Log.Log(LogLevelNotice, format, v...)
	}
}

func (me *LogModule) Info(module int, format string, v ...interface{}) {
	if me.isLog(module) {
		Log.Log(LogLevelInfo, format, v...)
	}
}

func (me *LogModule) Debug(module int, format string, v ...interface{}) {
	if me.isLog(module) {
		Log.Log(LogLevelDebug, format, v...)
	}
}
