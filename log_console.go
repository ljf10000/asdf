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
	fmt.Printf(format+Crlf, v...)
}

func (me *ConsoleLogger) Panic(format string, v ...interface{}) {
	fmt.Printf("Panic: "+format+Crlf, v...)

	go func() {
		panic(fmt.Sprintf("console-logger Panic: "+format+Crlf, v...))
	}()
}

func (me *ConsoleLogger) Emerg(format string, v ...interface{}) {
	if LogLevelEmerg <= me.level {
		fmt.Printf(format+Crlf, v...)
	}
}

func (me *ConsoleLogger) Alert(format string, v ...interface{}) {
	if LogLevelAlert <= me.level {
		fmt.Printf(format+Crlf, v...)
	}
}

func (me *ConsoleLogger) Crit(format string, v ...interface{}) {
	if LogLevelCrit <= me.level {
		fmt.Printf(format+Crlf, v...)
	}
}

func (me *ConsoleLogger) Error(format string, v ...interface{}) {
	if LogLevelError <= me.level {
		fmt.Printf(format+Crlf, v...)
	}
}

func (me *ConsoleLogger) Warning(format string, v ...interface{}) {
	if LogLevelWarning <= me.level {
		fmt.Printf(format+Crlf, v...)
	}
}

func (me *ConsoleLogger) Warn(format string, v ...interface{}) {
	if LogLevelWarn <= me.level {
		fmt.Printf(format+Crlf, v...)
	}
}

func (me *ConsoleLogger) Notice(format string, v ...interface{}) {
	if LogLevelNotice <= me.level {
		fmt.Printf(format+Crlf, v...)
	}
}

func (me *ConsoleLogger) Info(format string, v ...interface{}) {
	if LogLevelInfo <= me.level {
		fmt.Printf(format+Crlf, v...)
	}
}

func (me *ConsoleLogger) Debug(format string, v ...interface{}) {
	if LogLevelDebug <= me.level {
		fmt.Printf(format+Crlf, v...)
	}
}

func (me *ConsoleLogger) Close() error {
	return nil
}
