package asdf

import (
	"fmt"
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

var logTypes = EnumMapper{
	Enum: "LogType",
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

	*me = LogType(idx)

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

type LogLevel int

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
		Name:  "warning",
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

type LogMsg struct {
	s     string
	level LogLevel
}

func (me *LogMsg) String() string {
	return fmt.Sprintf("[%s] %s %s",
		me.level.Short(),
		time.Now().Format("2006-01-02@15:04:05"),
		me.s)
}

var Log ILogger = Console
