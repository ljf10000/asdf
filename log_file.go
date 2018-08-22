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
				s:     fmt.Sprintf(format, v...),
				level: level,
			}

			me.fd.WriteString(msg.String())
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
