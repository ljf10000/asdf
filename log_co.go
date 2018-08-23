package asdf

import (
	"fmt"
	"os"
	"time"
)

type CoLogChan chan LogMsg

func OpenCoLogger(prefix string, size int, show bool) error {
	ch := make(CoLogChan, size)

	logger := &CoLogger{
		level:  LogLevelDeft,
		prefix: prefix,
		ch:     ch,
	}

	Log = logger

	go logger.run(ch, show)

	return nil
}

type CoLogger struct {
	level  LogLevel
	prefix string
	fd     *os.File
	ch     CoLogChan
	year   int
	month  time.Month
	day    int
}

func (me *CoLogger) open() error {
	var err error

	file := fmt.Sprintf("%s.%04d-%02d-%02d.log",
		me.prefix,
		me.year,
		me.month,
		me.day)

	me.fd, err = os.OpenFile(file, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if nil != err {
		return err
	}

	return nil
}

func (me *CoLogger) tryCut() error {
	year, month, day := time.Now().Date()

	if year != me.year || month != me.month || day != me.day {
		me.Close()

		me.year, me.month, me.day = year, month, day

		return me.open()
	}

	return nil
}

func (me *CoLogger) run(ch CoLogChan, show bool) {
	fmt.Printf("cologger running...\n")

	for {
		msg, ok := <-ch
		if !ok {
			me.Close()

			return
		}

		s := msg.String()
		me.fd.WriteString(s)

		if show {
			fmt.Print(s)
		}

		err := me.tryCut()
		if nil != err {
			me.Close()

			return
		}
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
		me.ch <- LogMsg{
			s:     fmt.Sprintf(format, v...),
			level: level,
		}
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

func (me *CoLogger) Warn(format string, v ...interface{}) {
	me.Log(LogLevelWarn, format+Crlf, v...)
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
	// me.ch = nil

	err := me.fd.Close()

	me.fd = nil

	return err
}
