package asdf

import (
	"fmt"
	"os"
)

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
			msg := LogMsg{
				s:     fmt.Sprintf(format+Crlf, v...),
				level: level,
			}

			me.fd.WriteString(msg.String())
		})
	}
}

func (me *FileLogger) Panic(format string, v ...interface{}) {
	me.lock.Handle(func() {
		msg := LogMsg{
			s:     fmt.Sprintf("Panic: "+format+Crlf, v...),
			panic: true,
		}

		me.fd.WriteString(msg.String())
		me.fd.Close()
	})

	go func() {
		panic(fmt.Sprintf("file-logger Panic: "+format+Crlf, v...))
	}()
}

func (me *FileLogger) Emerg(format string, v ...interface{}) {
	if LogLevelEmerg <= me.level {
		me.Log(LogLevelEmerg, format, v...)
	}
}

func (me *FileLogger) Alert(format string, v ...interface{}) {
	if LogLevelAlert <= me.level {
		me.Log(LogLevelAlert, format, v...)
	}
}

func (me *FileLogger) Crit(format string, v ...interface{}) {
	if LogLevelCrit <= me.level {
		me.Log(LogLevelCrit, format, v...)
	}
}

func (me *FileLogger) Error(format string, v ...interface{}) {
	if LogLevelError <= me.level {
		me.Log(LogLevelError, format, v...)
	}
}

func (me *FileLogger) Warning(format string, v ...interface{}) {
	if LogLevelWarning <= me.level {
		me.Log(LogLevelWarning, format, v...)
	}
}

func (me *FileLogger) Warn(format string, v ...interface{}) {
	if LogLevelWarn <= me.level {
		me.Log(LogLevelWarn, format, v...)
	}
}

func (me *FileLogger) Notice(format string, v ...interface{}) {
	if LogLevelNotice <= me.level {
		me.Log(LogLevelNotice, format, v...)
	}
}

func (me *FileLogger) Info(format string, v ...interface{}) {
	if LogLevelInfo <= me.level {
		me.Log(LogLevelInfo, format, v...)
	}
}

func (me *FileLogger) Debug(format string, v ...interface{}) {
	if LogLevelDebug <= me.level {
		me.Log(LogLevelDebug, format, v...)
	}
}

func (me *FileLogger) Close() error {
	return me.fd.Close()
}
